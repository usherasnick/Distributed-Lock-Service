package distlock

import (
	"errors"
	"time"

	"github.com/FZambia/sentinel"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog/log"
)

type RedisConnPool struct {
	db int
	p  *redis.Pool
}

type RedisServiceConfig struct {
	SentinelEndpoints       []string `json:"sentinel_endpoints"`
	SentinelMasterName      string   `json:"sentinel_master_name"`
	SentinelPassword        string   `json:"sentinel_password"`
	RedisDatabase           int      `json:"redis_database"`
	RedisMasterPassword     string   `json:"redis_master_password"`
	RedisPoolMaxIdleConns   int      `json:"redis_pool_max_idle_conns"`   // 连接池最大空闲连接数
	RedisPoolMaxActiveConns int      `json:"redis_pool_max_active_conns"` // 连接池最大激活连接数
	RedisConnectTimeoutMsec int      `json:"redis_connect_timeout_msec"`  // 连接超时
	RedisReadTimeoutMsec    int      `json:"redis_read_timeout_msec"`     // 读取超时
	RedisWriteTimeoutMsec   int      `json:"redis_write_timeout_msec"`    // 写入超时
}

// EstablishRedisConn 建立连接redis服务的TCP连接池.
func EstablishRedisConn(cfg *RedisServiceConfig) *RedisConnPool {
	sntnl := &sentinel.Sentinel{
		Addrs:      cfg.SentinelEndpoints,
		MasterName: cfg.SentinelMasterName,
		Dial: func(addr string) (redis.Conn, error) {
			conn, err := redis.Dial(
				"tcp",
				addr,
				redis.DialConnectTimeout(time.Duration(cfg.RedisConnectTimeoutMsec)*time.Millisecond),
				redis.DialPassword(cfg.SentinelPassword),
			)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to connect to redis server")
				return nil, err
			}
			return conn, nil
		},
	}
	instance := &RedisConnPool{}
	instance.db = cfg.RedisDatabase
	instance.p = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			master, err := sntnl.MasterAddr()
			if err != nil {
				log.Fatal().Err(err).Msg("failed to connect to redis server")
				return nil, err
			}
			conn, err := redis.Dial(
				"tcp",
				master,
				redis.DialReadTimeout(time.Duration(cfg.RedisReadTimeoutMsec)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(cfg.RedisWriteTimeoutMsec)*time.Millisecond),
				redis.DialPassword(cfg.RedisMasterPassword),
			)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to connect to redis server")
				return nil, err
			}
			return conn, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if !sentinel.TestRole(conn, "master") {
				return errors.New("Role check failed")
			} else {
				return nil
			}
		},
		MaxIdle:   cfg.RedisPoolMaxIdleConns,
		MaxActive: cfg.RedisPoolMaxActiveConns,
		Wait:      true,
	}
	return instance
}

// CloseRedisConn 释放TCP连接池.
func CloseRedisConn(p *RedisConnPool) {
	if p != nil {
		p.p.Close()
	}
}

// ExecCommand 执行redis命令, 完成后自动归还连接.
func (p *RedisConnPool) ExecCommand(cmd string, args ...interface{}) (interface{}, error) {
	conn := p.getConn()
	defer conn.Close()
	return conn.Do(cmd, args...)
}

// ExecLuaScript 执行lua脚本, 完成后自动归还连接.
func (p *RedisConnPool) ExecLuaScript(src string, keyCount int, keysAndArgs ...interface{}) (interface{}, error) {
	conn := p.getConn()
	defer conn.Close()
	luaScript := redis.NewScript(keyCount, src)
	return luaScript.Do(conn, keysAndArgs...)
}

func (p *RedisConnPool) getConn() redis.Conn {
	conn := p.p.Get()
	conn.Do("SELECT", p.db) // nolint
	return conn
}
