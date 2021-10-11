package cache

import (
	"fmt"
)

// Instance The type of function that created the instance
type Instance func() Cache

// Initialize cache driven map
var drivers = make(map[string]Instance)

// Cache Interface definition
type Cache interface {
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

	// StartInstance Start Cache instance
	StartInstance(config string) error
}

// Register Driver registration, the same driver registration twice will panic
func Register(name string, driverName Instance) {
	if driverName == nil {
		panic("cache: Register driver is nil")
	}
	if _, ok := drivers[name]; ok {
		panic("cache: This driver is already registered " + name)
	}
	drivers[name] = driverName
}

// New Create a new driver and pass in the configuration item
// The configuration item is in JSON format
func New(adapterName, config string) (adapter Cache, err error) {
	instanceFunc, ok := drivers[adapterName]
	if !ok {
		err = fmt.Errorf("cache: unknown driver name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.StartInstance(config)
	if err != nil {
		adapter = nil
	}
	return
}
