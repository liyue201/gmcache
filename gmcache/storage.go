package gmcache

import (
	"github.com/codinl/go-logger"
	"math/rand"
	"sync"
	"time"
)

type KVItem struct {
	value      []byte
	expiration time.Duration
	keyIdx     int //key index int Storage.keys
}

type Storage struct {
	sync.RWMutex
	keys           []string           //keys slide
	m              map[string]*KVItem //key to item
	memChangedChan chan int64
}

func (this *KVItem) expired() bool {
	if int64(this.expiration) < time.Now().UnixNano() {
		return true
	}
	return false
}

func NewStorage(memChangedChan chan int64) *Storage {
	return &Storage{
		keys:           make([]string, 0),
		m:              make(map[string]*KVItem),
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
		item.expiration = time.Duration(time.Now().UnixNano()) + ttl
	} else {
		newItem := &KVItem{
			value:      value,
			expiration: time.Duration(time.Now().UnixNano()) + ttl,
			keyIdx:     len(this.keys),
		}
		this.m[key] = newItem
		this.keys = append(this.keys, key)
		this.memUsedChanged(int64(len(key)*2 + len(value)))
	}

	return nil
}

func (this *Storage) Get(key string) (*KVItem, error) {
	this.RLock()
	item, ok := this.m[key]
	this.RUnlock()

	if ok {
		if item.expired() {

			this.Lock()  //Prevent duplicate delete
			item, ok := this.m[key]
			if ok {
				this.deleteItem(key, item)
				this.memUsedChanged(int64(-len(key)*2 - len(item.value)))
			}
			this.Unlock()

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
		this.deleteItem(key, item)
		this.memUsedChanged(int64(-len(key)*2 - len(item.value)))
	}
	return nil
}

func (this *Storage) deleteItem(key string, item *KVItem) {
	lastKey := this.keys[len(this.keys)-1]
	this.keys[item.keyIdx] = lastKey
	this.m[lastKey].keyIdx = item.keyIdx

	this.keys = this.keys[:len(this.keys)-1]
	delete(this.m, key)

	logger.Debug("delete:", key)
}

func (this *Storage) itemNum() int {
	this.Lock()
	defer this.Unlock()
	return len(this.keys)
}

func (this *Storage) memUsedChanged(bytes int64) {
	this.memChangedChan <- bytes
}

//Refer to redis's mechanism
//1.Test 100 keys
//2.Delete all expired keys
//3.Repet step 1 if there are over 25 keys have been deleted

func (this *Storage) DeleteExpiredKeyRandom() int64 {
	now := time.Now()
	r := rand.New(rand.NewSource(now.Unix()))

	var totalDelBytes int64

	for {
		itemNum := this.itemNum()
		deletedCount := 0
		for i := 0; i < 100 && i < itemNum; i++ {
			this.Lock()
			n := len(this.keys)
			if n == 0 {
				this.Unlock()
				break
			}
			index := int(r.Int31n(int32(n)))
			key := this.keys[index]
			item := this.m[key]
			if int64(item.expiration) < now.UnixNano() {
				this.deleteItem(key, item)
				totalDelBytes += int64(len(key) + len(item.value))
				deletedCount++
			}
			this.Unlock()
		}

		if deletedCount < 25 {
			break
		}
	}
	return totalDelBytes
}
