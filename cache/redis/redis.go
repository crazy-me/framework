package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crazy-me/framework/cache"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

// Cache Interface Implementation
type Cache struct {
	c        *redis.Pool   // Redis connection pool
	address  string        // Redis address 127.0.0.1:6379
	password string        // Redis password
	db       int           // Redis database number
	prefix   string        // Redis key prefix
	maxIdle  int           // Pool Maximum number of free connections
	timeout  time.Duration // Close connections after remaining idle for this duration
}

// NewRedis Create a new redis
func NewRedis() cache.Cache {
	return &Cache{}
}

// Get Cache Value
func (rc *Cache) Get(key string) interface{} {
	if v, err := rc.do("GET", key); err == nil {
		return v
	}
	return nil
}

// Set Cache Value
func (rc *Cache) Set(key string, val interface{}) bool {
	if _, err := rc.do("SET", key, val); err == nil {
		return true
	}
	return false
}

// SetEx Cache Value By timeout
func (rc *Cache) SetEx(key string, val interface{}, expiration int64) bool {
	if _, err := rc.do("SET", key, val, "EX", expiration); err == nil {
		return true
	}
	return false
}

// GetMulti Get batch Cache Value
func (rc *Cache) GetMulti(keys []string) []interface{} {
	r := rc.c.Get()
	defer r.Close()
	var args []interface{}
	for _, key := range keys {
		args = append(args, rc.associate(key))
	}
	values, err := redis.Values(r.Do("MGET", args...))
	if err != nil {
		return nil
	}
	return values
}

// Delete Cache Value By Key
func (rc *Cache) Delete(key string) bool {
	if reply, _ := rc.do("DEL", key); reply == int64(0) {
		return false
	}
	return true
}

// Incr Increase Cache Value
func (rc *Cache) Incr(key string) bool {
	if _, err := rc.do("INCR", key); err == nil {
		return true
	}
	return false
}

// Decr Reduce Cache Value
func (rc *Cache) Decr(key string) bool {
	if _, err := rc.do("DECR", key); err == nil {
		return true
	}
	return false
}

// IsExist Check Key Does it exist
func (rc *Cache) IsExist(key string) bool {
	if reply, _ := rc.do("EXISTS", key); reply == int64(0) {
		return false
	}
	return true
}

// ClearAll clear all Cache Value By Key
func (rc *Cache) ClearAll() error {
	r := rc.c.Get()
	defer r.Close()
	_, err := r.Do("FLUSHDB")
	return err
}

// do actually do the redis cmd
func (rc *Cache) do(command string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	args[0] = rc.associate(args[0])
	r := rc.c.Get()
	defer r.Close()
	return r.Do(command, args...)
}

// associate with config key.
func (rc *Cache) associate(originKey interface{}) string {
	return fmt.Sprintf("%s:%s", rc.prefix, originKey)
}

// StartInstance redis
// {"address":":6379","password":"","db":"0","prefix":"","maxIdle":"","timeout":""}
func (rc *Cache) StartInstance(config string) error {
	var cf map[string]string
	_ = json.Unmarshal([]byte(config), &cf)

	// Instance configuration
	if _, ok := cf["address"]; !ok {
		return errors.New("no address configure")
	}

	if _, ok := cf["password"]; !ok {
		cf["password"] = ""
	}

	if _, ok := cf["db"]; !ok {
		cf["db"] = "0"
	}

	if _, ok := cf["prefix"]; !ok {
		cf["prefix"] = "Redis"
	}

	if _, ok := cf["maxIdle"]; !ok {
		cf["maxIdle"] = "3"
	}

	if _, ok := cf["timeout"]; !ok {
		cf["timeout"] = "180s"
	}

	rc.prefix = cf["prefix"]
	rc.address = cf["address"]
	rc.db, _ = strconv.Atoi(cf["db"])
	rc.password = cf["password"]
	rc.maxIdle, _ = strconv.Atoi(cf["maxIdle"])
	if v, err := time.ParseDuration(cf["timeout"]); err == nil {
		rc.timeout = v
	} else {
		rc.timeout = 180 * time.Second
	}

	// Redis Pool Connection
	rc.connection()

	r := rc.c.Get()
	defer r.Close()

	return r.Err()
}

func (rc *Cache) connection() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rc.address)
		if err != nil {
			return nil, err
		}

		if rc.password != "" {
			if _, err := c.Do("AUTH", rc.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, base := c.Do("SELECT", rc.db)
		if base != nil {
			c.Close()
			return nil, base
		}
		return
	}

	// initialize a new pool
	rc.c = &redis.Pool{
		MaxIdle:     rc.maxIdle,
		IdleTimeout: rc.timeout,
		Dial:        dialFunc,
	}
}

// Register redis driver
func init() {
	cache.Register("redis", NewRedis)
}
