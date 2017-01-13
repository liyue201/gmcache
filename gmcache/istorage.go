package gmcache

import (
	"errors"
	"time"
)

var EXPIRED_ERROR error = errors.New("Expire error")
var KEY_NOT_EXIST_ERROR error = errors.New("Expire error")
var OUT_OF_MEMORY_LIMIT_ERROR error = errors.New("Out of memory limit error")

type IStorage interface {
	Set(key string, value []byte, ttl time.Duration) error
	Get(key string) (*KVItem, error)
	Delete(key string) error
}
