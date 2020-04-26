package distlock

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

const (
	_DistributedLock = "pddlock"

	/*
		-1: failed to get
		 0: failed to del
		 1: success to del
	*/
	_CheckAndDel = `if redis.call('get', KEYS[1]) == ARGV[1] then
	return redis.call('del', KEYS[1])
else
	return -1
end`
)

// DLockByRedis 通过redis实现的分布式锁
type DLockByRedis struct {
	p *RedisConnPool
}

// NewDLockByRedis 获取DLockByRedis实例.
func NewDLockByRedis(p *RedisConnPool) *DLockByRedis {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Caller().Logger()

	return &DLockByRedis{
		p: p,
	}
}

// TryLock 尝试获取分布式锁, 超时后就放弃 (不可重入锁).
func (dlr *DLockByRedis) TryLock(timeoutInSecs int64) (string, bool) {
	id := uuid.NewV1().String()
	ttr := time.Now().Unix() + timeoutInSecs
	for {
		// 为了避免出现死锁状态, 需要设置一个合理的过期时间
		// TODO: 设置为多少比较合理?
		v, err := dlr.p.ExecCommand("SET", _DistributedLock, id, "NX", "EX", 3600)
		if err != nil {
			log.Error().Err(err).Msg("failed to acquire lock")
			return "", false
		}
		if v == nil {
			continue
		}
		if v.(string) == "OK" {
			return id, true
		}
		if time.Now().Unix() > ttr {
			return "", false
		}
	}
}

// Unlock 释放分布式锁.
func (dlr *DLockByRedis) Unlock(value string) {
	v, err := dlr.p.ExecLuaScript(_CheckAndDel, 1, _DistributedLock, value)
	if err != nil {
		log.Error().Err(err).Msg("failed to release lock")
	}
	if v == nil || v.(int64) != 1 {
		log.Error().Err(err).Msg("failed to release lock")
	}
}
