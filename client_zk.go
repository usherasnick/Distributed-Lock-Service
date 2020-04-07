package distlock

import (
	"bytes"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	_LockerRootPath                           = "/distlock"
	_LockerLockPathFastLockUsed               = "/distlock/fast-lock"
	_LockerLockPathFastLockUsedPrefix         = "/distlock/fast-lock/request-"
	_LockerLockPathFastLockUsedShortestPrefix = "request-"

	_MaxRetries = 3
)

// EstablishZKConn 建立一条连接zookeeper集群的TCP连接.
func EstablishZKConn(endpoints []string) (*zk.Conn, <-chan zk.Event) {
	conn, evCh, err := zk.Connect(endpoints, time.Second*time.Duration(10))
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to connect to zookeeper cluster (%v)", endpoints)
	}
	createIfNotExistOrDie(conn, _LockerRootPath)
	createIfNotExistOrDie(conn, _LockerLockPathFastLockUsed)
	return conn, evCh
}

// CloseZKConn 关闭TCP连接.
func CloseZKConn(conn *zk.Conn) {
	conn.Close()
}

// 如果ZNode不存在就创建
func createIfNotExist(conn *zk.Conn, path string) error {
	if _, err := safeCreate(conn, path, []byte(""), 0); err != nil && err != zk.ErrNodeExists {
		return err
	}
	return nil
}

// 如果ZNode不存在就创建, 出错直接panic
func createIfNotExistOrDie(conn *zk.Conn, path string) {
	if err := createIfNotExist(conn, path); err != nil {
		log.Fatal().Err(err).Msgf("failed to create znode <%s>", path)
	}
}

// 创建ZNode
func safeCreate(conn *zk.Conn, path string, data []byte, flags int32) (string, error) {
	var _path string
	var err error
	var _data []byte

	var retry bool

LOOP:
	for i := 0; i < _MaxRetries; i++ {
		_path, err = conn.Create(path, data, flags, zk.WorldACL(zk.PermAll))
		switch err {
		// No need to search for the node since it can't exist. Just try again.
		case zk.ErrSessionExpired:
			{
				continue
			}
		// 连接关闭, 可能因为暂时的网络问题, 直接重试
		case zk.ErrConnectionClosed:
			{
				retry = true
				continue
			}
		// ZNode已存在
		case zk.ErrNodeExists:
			{
				// 之前就创建过
				if !retry {
					return _path, zk.ErrNodeExists
				}
				// 因为网络问题导致的假失败
				_data, _, err = safeGet(conn, path)
				if err != nil {
					// 又可能因为暂时的网络问题, 请重试
					continue
				}
				if bytes.Equal(data, _data) {
					return _path, nil
				}
				return "", zk.ErrUnknown
			}
		// TODO: 处理更多的错误情形
		default:
			{
				break LOOP
			}
		}
	}

	return _path, err
}

// 获取ZNode的值
func safeGet(conn *zk.Conn, path string) ([]byte, *zk.Stat, error) {
	var _data []byte
	var _stat *zk.Stat
	var err error

LOOP:
	for i := 0; i < _MaxRetries; i++ {
		_data, _stat, err = conn.Get(path)
		switch err {
		// session过期直接panic
		case zk.ErrSessionExpired:
			{
				log.Fatal().Err(zk.ErrSessionExpired).Msgf("failed to get value of znode <%s>", path)
			}
		// 连接关闭, 可能因为暂时的网络问题, 请重试
		case zk.ErrConnectionClosed:
			{
				continue
			}
		// TODO: 处理更多的错误情形
		default:
			{
				break LOOP
			}
		}
	}

	return _data, _stat, err
}

// 删除ZNode
func safeDelete(conn *zk.Conn, path string, version int32) error {
	var err error
	var retry bool

LOOP:
	for i := 0; i < _MaxRetries; i++ {
		err = conn.Delete(path, version)
		switch err {
		// session过期直接panic
		case zk.ErrSessionExpired:
			{
				log.Fatal().Err(zk.ErrSessionExpired).Msgf("failed to delete znode <%s>", path)
			}
		// 连接关闭, 可能因为暂时的网络问题, 请重试
		case zk.ErrConnectionClosed:
			{
				retry = true
				continue
			}
		// ZNode不存在
		case zk.ErrNoNode:
			{
				// 因为网络问题导致的假失败
				if retry {
					return nil
				}
				return zk.ErrNoNode
			}
		// TODO: 处理更多的错误情形
		default:
			{
				break LOOP
			}
		}
	}

	return err
}

// 获取ZNode所有的子节点
func safeGetChildren(conn *zk.Conn, path string, watch bool) ([]string, <-chan zk.Event, error) {
	var _children []string
	var _watcher <-chan zk.Event
	var err error

LOOP:
	for i := 0; i < _MaxRetries; i++ {
		if watch {
			_children, _, _watcher, err = conn.ChildrenW(path)
		} else {
			_children, _, err = conn.Children(path)
		}
		switch err {
		// session过期直接panic
		case zk.ErrSessionExpired:
			{
				log.Fatal().Err(zk.ErrSessionExpired).Msgf("failed to list children of znode <%s>", path)
			}
		// 连接关闭, 可能因为暂时的网络问题, 请重试
		case zk.ErrConnectionClosed:
			{
				continue
			}
		// TODO: 处理更多的错误情形
		default:
			{
				break LOOP
			}
		}
	}

	return _children, _watcher, err
}
