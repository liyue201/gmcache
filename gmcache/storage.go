package gmcache

import (
	"math/rand"
	"sync"
	"time"
)

//the implementation of IStorage
type Storage struct {
	sync.RWMutex
	m              map[string]*KVItem
	keys           []string
	keyIdx         map[string]int //key index in keys
	memChangedChan chan int64
}

func NewStorage(memChangedChan chan int64) *Storage {
	return &Storage{
		m:              make(map[string]*KVItem),
		keys:           make([]string, 0),
		keyIdx:         make(map[string]int),
		memChangedChan: memChangedChan,
	}
}

func (this *Storage) Set(key string, value []byte, ttl time.Duration) error {
	this.Lock()
	defer this.Unlock()

	item, ok := this.m[key]
	if ok {
		this.memUsedChanged(int64(len(value) - len(item.value)))
		item.value = value
		item.expire = time.Duration(time.Now().UnixNano()) + ttl
	} else {
		newItem := &KVItem{
			key:    key,
			value:  value,
			expire: time.Duration(time.Now().UnixNano()) + ttl,
		}
		this.m[key] = newItem
		this.keys = append(this.keys, key)
		this.keyIdx[key] = len(this.keys) - 1
		this.memUsedChanged(int64(len(key) + len(value)))
	}

	return nil
}

func (this *Storage) Get(key string) (*KVItem, error) {
	this.RLock()
	defer this.RUnlock()

	item, ok := this.m[key]
	if ok {
		if item.expired() {
			this.deleteItem(item)
			this.memUsedChanged(int64(-len(item.key) - len(item.value)))
			return nil, EXPIRED_ERROR
		}
		return item, nil
	}
	return nil, KEY_NOT_EXIST_ERROR
}

func (this *Storage) Delete(key string) error {
	this.Lock()
	defer this.Unlock()

	item, ok := this.m[key]
	if ok {
		this.deleteItem(item)
		this.memUsedChanged(int64(-len(item.key) - len(item.value)))
	}
	return nil
}

func (this *Storage) deleteItem(item *KVItem) {
	delete(this.m, item.key)
	this.keys[this.keyIdx[item.key]] = this.keys[len(this.keys)-1]
	this.keys = this.keys[:len(this.keys)-1]
	delete(this.keyIdx, item.key)
}

func (this *Storage) KeysNum() int {
	this.Lock()
	defer this.Unlock()
	return len(this.keys)
}

func (this *Storage) memUsedChanged(bytes int64) {
	this.memChangedChan <- bytes
}

//Refer to redis's mechanism
//1.Rest 100 keys
//2.Delete all expired keys
//3.Repet step 1 if there are over 25 keys have been deleted

func (this *Storage) DeleteExpiredKeyRandom() int64 {
	now := time.Now()
	r := rand.New(rand.NewSource(now.Unix()))

	var totalDelBytes int64

	for {
		if this.KeysNum() < 100 {
			break
		}

		deletedCount := 0
		for i := 0; i < 100; i++ {
			this.Lock()
			n := len(this.keys)
			if n == 0 {
				this.Unlock()
				break
			}
			index := int(r.Int31n(int32(n)))
			key := this.keys[index]
			item := this.m[key]
			if int64(item.expire) < now.UnixNano() {
				this.deleteItem(item)
				totalDelBytes += int64(len(item.key) + len(item.value))
				deletedCount++
			}
			this.Unlock()
		}

		if deletedCount < 25 {
			break
		}
	}
	return  totalDelBytes
}
