## cache install

	go get github.com/crazy-me/framework/cache

## cache init

    import (
		_"github.com/crazy-me/framework/cache/redis"
	)

## cache use

	import (
		"github.com/crazy-me/framework/cache"
	)

## init a cache
    bm, err := cache.New("redis", `{"address":":6379","password":"","db":"0","maxIdle":"3","timeout":""}`)

## usage method

    // Get Cache Value
	Get(key string) interface{}

	// Set Cache Value
	Set(key string, val interface{}) bool

	// SetEx Cache Value By timeout
	SetEx(key string, val interface{}, expiration int64) bool

	// GetMulti Get batch Cache Value
	GetMulti(keys []string) []interface{}

	// Delete Cache Value By Key
	Delete(key string) bool

	//Incr Increase Cache Value
	Incr(key string) bool

	// Decr Reduce Cache Value
	Decr(key string) bool

	// IsExist Check Key Does it exist
	IsExist(key string) bool

	// ClearAll clear all Cache Value By Key
	ClearAll() error
