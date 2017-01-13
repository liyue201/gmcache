package gmcache

import (
	"sync"
	"time"
)

//the implementation of IStorage
type Storage struct {
	sync.RWMutex
	m map[string]*KVItem
}

func NewStorage() IStorage {
	return &Storage{
		m: make(map[string]*KVItem),
	}
}

func (this *Storage) Set(key string, value []byte, ttl time.Duration) error {
	item := &KVItem{
		key:    key,
		value:  value,
		expire: time.Duration(time.Now().UnixNano()) + ttl,
	}

	this.Lock()
	defer this.Unlock()
	this.m[key] = item
	return nil
}

func (this *Storage) Get(key string) (*KVItem, error) {
	this.RLock()
	defer this.RUnlock()

	item, ok := this.m[key]
	if ok {
		if item.expired() {
			return nil, EXPIRED_ERROR
		}
		return item, nil
	}
	return nil, KEY_NOT_EXIST_ERROR
}

func (this *Storage) Delete(key string) error {
	this.Lock()
	defer this.Unlock()

	_, ok := this.m[key]
	if ok {
		delete(this.m, key)
	}
	return nil
}
