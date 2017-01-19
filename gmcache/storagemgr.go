package gmcache

import (
	"hash/crc32"
	"sync/atomic"
	"time"
)

//The implementation of IStorage
type StorageManager struct {
	bucketNum      int
	buckets        []*Storage
	memoryLimit    int64
	memeryUsed     int64
	cleanInterval  time.Duration
	memChangedChan chan int64
	stop           chan struct{}
}

//memoryLimit unit is in byte
func NewStorageManager(bucketNum int, memoryLimit int64, cleanInterval time.Duration) *StorageManager {
	memMgr := &StorageManager{
		bucketNum:      bucketNum,
		buckets:        make([]*Storage, bucketNum),
		memoryLimit:    memoryLimit,
		memeryUsed:     0,
		cleanInterval:  cleanInterval,
		memChangedChan: make(chan int64, 1000),
		stop:           make(chan struct{}),
	}
	for i := 0; i < bucketNum; i++ {
		memMgr.buckets[i] = NewStorage(memMgr.memChangedChan)
	}
	return memMgr
}

func (this *StorageManager) mapToIndex(key string) int {
	return int(crc32.ChecksumIEEE([]byte(key))) % this.bucketNum
}

func (this *StorageManager) findStorage(key string) *Storage {
	return this.buckets[this.mapToIndex(key)]
}

func (this *StorageManager) Set(key string, value []byte, ttl time.Duration) error {
	if (int64(len(value)) + this.memeryUsed) > this.memoryLimit {
		return OUT_OF_MEMORY_LIMIT_ERROR
	}
	return this.findStorage(key).Set(key, value, ttl)
}

func (this *StorageManager) Get(key string) (*KVItem, error) {
	return this.findStorage(key).Get(key)
}

func (this *StorageManager) Delete(key string) error {
	return this.findStorage(key).Delete(key)
}

func (this *StorageManager) Run() {
	cleanTicker := time.NewTicker(this.cleanInterval)

	for {
		select {
		case <-cleanTicker.C:
			for i := 0; i < this.bucketNum; i++ {
				deletedBytes := this.buckets[i].DeleteExpiredKeyRandom()
				atomic.AddInt64(&this.memeryUsed, -deletedBytes)
			}
		case deltaBytes := <-this.memChangedChan:
			atomic.AddInt64(&this.memeryUsed, deltaBytes)
		case <-this.stop:
			return
		}
	}
}

func (this *StorageManager) Stop() {
	close(this.stop)
}
