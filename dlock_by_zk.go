package distlock

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samuel/go-zookeeper/zk"
)

// DLockByZookeeper 通过zookeeper实现的分布式锁
type DLockByZookeeper struct {
	conn     *zk.Conn
	lockpath string
}

// NewDLockByZookeeper 获取DLockByZookeeper实例.
func NewDLockByZookeeper(conn *zk.Conn) *DLockByZookeeper {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Caller().Logger()

	return &DLockByZookeeper{
		conn: conn,
	}
}

// TryLock 尝试获取分布式锁, 超时后就放弃 (不可重入锁).
/*

==> acquire lock
n = create("/distlock/fast-lock/request-", "", ephemeral|sequence)
RETRY:
    children = getChildren("/distlock/fast-lock", watch=False)
    if n is lowest znode in children:
        return
    else:
        exist("/distlock/fast-lock/request-" % (n - 1), watch=True)

watch_event:
	goto RETRY

*/
func (dlz *DLockByZookeeper) TryLock(timeoutInSecs int) bool {
	if dlz.lockpath != "" {
		log.Error().Err(ErrDeadlock).Msg("failed to acquire lock")
		return false
	}

	path, err := safeCreate(dlz.conn, _LockerLockPathFastLockUsedPrefix, []byte(""), zk.FlagEphemeral|zk.FlagSequence)
	if err != nil {
		log.Error().Err(err).Msg("failed to acquire lock")
		return false
	}
	seq := dlz.getSequenceNum(path, _LockerLockPathFastLockUsedPrefix)

	ticker := time.NewTicker(time.Duration(timeoutInSecs) * time.Second)
LOOP:
	for {
		select {
		case <-ticker.C:
			{
				log.Warn().Msg("failed to acquire lock, since timeout to retry")
				return false
			}
		default:
			{
			TRY_AGAIN:
				children, _, err := safeGetChildren(dlz.conn, _LockerLockPathFastLockUsed, false)
				if err != nil {
					log.Error().Err(err).Msg("failed to acquire lock")
					return false
				}

				minSeq := seq
				prevSeq := -1
				prevSeqPath := ""

				for _, child := range children {
					s := dlz.getSequenceNum(child, _LockerLockPathFastLockUsedShortestPrefix)
					if s < minSeq {
						minSeq = s
					}
					if s < seq && s > prevSeq {
						prevSeq = s
						prevSeqPath = child
					}
				}

				if seq == minSeq {
					break LOOP
				}

				_, _, watcher, err := dlz.conn.ExistsW(_LockerLockPathFastLockUsed + "/" + prevSeqPath)
				if err != nil {
					log.Error().Err(err).Msg("failed to acquire lock")
					return false
				}

				smallTicker := time.NewTicker(time.Duration(timeoutInSecs/3) * time.Second)
				for {
					select {
					case ev, ok := <-watcher:
						{
							if !ok {
								return false
							}
							if ev.Type == zk.EventNodeDeleted {
								goto TRY_AGAIN
							}
						}
					case <-smallTicker.C:
						{
							goto TRY_AGAIN
						}
					}
				}
			}
		}
	}

	dlz.lockpath = path
	return true
}

// Unlock 释放分布式锁.
/*

==> release lock (voluntarily or session timeout)
delete("/distlock/fast-lock/request-" % n)

*/
func (dlz *DLockByZookeeper) Unlock() {
	if err := safeDelete(dlz.conn, dlz.lockpath, -1); err != nil {
		log.Error().Err(err).Msg("failed to release lock")
	}
	dlz.lockpath = ""
}

func (dlz *DLockByZookeeper) getSequenceNum(path, prefix string) int {
	numStr := strings.TrimPrefix(path, prefix)
	num, _ := strconv.Atoi(numStr)
	return num
}
