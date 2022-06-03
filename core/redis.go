package core

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// Cache  for caching
type Cache struct {
	Pool *redis.Pool
}

// CacheService for caching service
type CacheService interface {
	Ping() error
	Get(key string) ([]byte, error)
	Set(key string, value []byte, ttl int64) error
	Exists(key string) (bool, error)
	Delete(key string) error
	Increment(key string) (int64, error)
	IncrementWithTTL(key string, ttl int64) (int64, error)
	IncrementWithValue(key string, value float64) (int64, error)
	IncrementWithValueTTL(key string, ttl int64, value float64) (int64, error)
}

var cache CacheService

// NewCacheService for new cache service
func NewCacheService(pool *redis.Pool) CacheService {
	if cache == nil {
		cache = &Cache{pool}
	}
	return cache
}

// CreateRedisPool for connection pool redis
func CreateRedisPool(addr, password string, maxIdle int) (*redis.Pool, error) {
	redis := &redis.Pool{
		MaxIdle: maxIdle,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}

			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					return nil, err
				}
			}
			return c, nil
		},
	}

	conn := redis.Get()
	_, err := conn.Do("PING")
	if err != nil {
		return redis, err
	}

	defer conn.Close()

	return redis, nil
}

// Ping ping a server
func (cache *Cache) Ping() error {

	conn := cache.Pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		return fmt.Errorf("cannot 'PING' db: %v", err)
	}
	return nil
}

// Get get value from key
func (cache *Cache) Get(key string) ([]byte, error) {

	conn := cache.Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("error getting key %s: %v", key, err)
	}
	return data, err
}

// Set set key, value
func (cache *Cache) Set(key string, value []byte, ttl int64) error {

	conn := cache.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	_, err = conn.Do("EXPIRE", key, ttl)

	if err != nil {
		return fmt.Errorf("error setting expire key %s to %d: %v", key, ttl, err)
	}
	return nil

}

// Exists check key is exist
func (cache *Cache) Exists(key string) (bool, error) {

	conn := cache.Pool.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}
	return ok, err
}

// Increment check key is exist
func (cache *Cache) Increment(key string) (int64, error) {

	conn := cache.Pool.Get()
	defer conn.Close()

	ok, err := redis.Int64(conn.Do("INCR", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}

	return ok, err
}

// IncrementWithTTL for increment with counter
func (cache *Cache) IncrementWithTTL(key string, ttl int64) (int64, error) {
	val, err := cache.Increment(key)
	if err != nil {
		return val, err
	}

	conn := cache.Pool.Get()
	_, err = conn.Do("EXPIRE", key, ttl)

	if err != nil {
		return val, fmt.Errorf("error setting expire key %s to %d: %v", key, ttl, err)
	}
	return val, nil
}

// IncrementWithValue for increment with counter
func (cache *Cache) IncrementWithValue(key string, value float64) (int64, error) {
	conn := cache.Pool.Get()
	defer conn.Close()

	v, err := redis.Int64(conn.Do("INCRBY", key, value))
	if err != nil {
		return v, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}

	return v, err
}

// IncrementWithValueTTL for increment with value with ttl
func (cache *Cache) IncrementWithValueTTL(key string, ttl int64, value float64) (int64, error) {
	conn := cache.Pool.Get()
	defer conn.Close()

	v, err := redis.Int64(conn.Do("INCRBY", key, value))
	if err != nil {
		return v, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}

	_, err = conn.Do("EXPIRE", key, ttl)
	if err != nil {
		return v, fmt.Errorf("error setting expire key %s to %d: %v", key, ttl, err)
	}

	return v, err
}

// Delete delete by keys
func (cache *Cache) Delete(key string) error {

	conn := cache.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}
