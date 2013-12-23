package caching

import (
    "time"
    "github.com/robfig/revel"
    "github.com/mawenbao/wishome/app"
    "github.com/garyburd/redigo/redis"
)

// global redis connection pool
var RedisConnPool *redis.Pool

func init() {
    revel.OnAppStart(initRedisConnPool)
}

func initRedisConnPool() {
    RedisConnPool = &redis.Pool {
        MaxIdle: app.MyGlobal.RedisPoolMaxIdle,
        IdleTimeout: app.MyGlobal.RedisPoolIdleTimeout,
        MaxActive: 0, // no active connection limit
        Dial: newRedisConn,
        TestOnBorrow: testRedisConn,
    }
}

// no password
func newRedisConn() (redis.Conn, error) {
    c, err := redis.Dial("tcp", app.MyGlobal.RedisServerAddr)
    if nil != err {
        return nil, err
    }
    return c, nil
}

func testRedisConn(c redis.Conn, t time.Time) error {
    _, err := c.Do("PING")
    return err
}

func RedisConn() redis.Conn {
    return RedisConnPool.Get()
}

