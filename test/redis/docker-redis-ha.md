## 搭建Redis HA架构

1. 创建HA架构可用的网络环境

```shell
docker network create redis-ha --subnet=192.168.1.0/24
```

2. 先启动主节点

```shell
docker run --name=redis_master \
    --net=redis-ha \
    --ip=192.168.1.2 \
    -v $PWD/conf.d/redis-master.conf:/usr/local/etc/redis/redis.conf \
    -v ~/redis/data00:/data \
    -d redis:6.2-rc-alpine3.12 sh -c "redis-server /usr/local/etc/redis/redis.conf"
```

3. 再启动2个从节点

```shell
docker run --name=redis_slave01 \
    --net=redis-ha \
    --ip=192.168.1.3 \
    -v $PWD/conf.d/redis-slave01.conf:/usr/local/etc/redis/redis.conf \
    -v ~/redis/data01:/data \
    -d redis:6.2-rc-alpine3.12 sh -c "redis-server /usr/local/etc/redis/redis.conf"
```

```shell
docker run --name=redis_slave02 \
    --net=redis-ha \
    --ip=192.168.1.4 \
    -v $PWD/conf.d/redis-slave02.conf:/usr/local/etc/redis/redis.conf \
    -v ~/redis/data02:/data \
    -d redis:6.2-rc-alpine3.12 sh -c "redis-server /usr/local/etc/redis/redis.conf"
```

4. 再启动三个哨兵节点

```shell
docker run --name=redis_sentinel01 \
    --net=redis-ha \
    --ip=192.168.1.5 \
    -p 26379:26379 \
    -v $PWD/conf.d/sentinel.conf:/usr/local/etc/redis/sentinel.conf \
    -v ~/redis/sentinel/data00:/data \
    -d redis:6.2-rc-alpine3.12 sh -c "redis-sentinel /usr/local/etc/redis/sentinel.conf"
```

```shell
docker run --name=redis_sentinel02 \
    --net=redis-ha \
    --ip=192.168.1.6 \
    -p 26380:26379 \
    -v $PWD/conf.d/sentinel.conf:/usr/local/etc/redis/sentinel.conf \
    -v ~/redis/sentinel/data01:/data \
    -d redis:6.2-rc-alpine3.12 sh -c "redis-sentinel /usr/local/etc/redis/sentinel.conf"
```

```shell
docker run --name=redis_sentinel03 \
    --net=redis-ha \
    --ip=192.168.1.7 \
    -p 26381:26379 \
    -v $PWD/conf.d/sentinel.conf:/usr/local/etc/redis/sentinel.conf \
    -v ~/redis/sentinel/data02:/data \
    -d redis:6.2-rc-alpine3.12 sh -c "redis-sentinel /usr/local/etc/redis/sentinel.conf"
```
