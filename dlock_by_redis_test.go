package distlock

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	fakeRedisServiceCfg = &RedisServiceConfig{
		SentinelEndpoints:       []string{"localhost:26379", "localhost:26380", "localhost:26381"},
		SentinelMasterName:      "mymaster",
		SentinelPassword:        "Pwd123!@",
		RedisMasterPassword:     "sOmE_sEcUrE_pAsS",
		RedisPoolMaxIdleConns:   3,
		RedisPoolMaxActiveConns: 64,
		RedisConnectTimeoutMsec: 500,
		RedisReadTimeoutMsec:    500,
		RedisWriteTimeoutMsec:   500,
	}
)

func TestDLockByRedis(t *testing.T) {
	p := EstablishRedisConn(fakeRedisServiceCfg)
	defer CloseRedisConn(p)

	total := 0

	var n sync.WaitGroup
	for i := 0; i < 200; i++ {
		n.Add(1)
		go func(p *RedisConnPool, idx int) {
			defer n.Done()

			dl := NewDLockByRedis(p)
			if v, ok := dl.TryLock(5); ok && v != "" {
				total++
				time.Sleep(time.Microsecond * 10)
				dl.Unlock(v)
			}
		}(p, i)
	}
	n.Wait()

	assert.Equal(t, 200, total)
}
