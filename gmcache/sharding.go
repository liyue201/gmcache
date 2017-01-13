package gmcache

import (
	"hash/crc32"
	"time"
)

//The implementation of IStorage
type ShardingStorage struct {
	bucketNum int
	buckets   []IStorage
}

func NewSharingStorage(bucketNum int) IStorage {
	shaldingStorage := &ShardingStorage{
		bucketNum: bucketNum,
		buckets:   make([]IStorage, bucketNum),
	}

	for i := 0; i < bucketNum; i++ {
		shaldingStorage.buckets[i] = NewStorage()
	}
	return shaldingStorage
}

func (this *ShardingStorage) mapToIndex(key string) int {
	return int(crc32.ChecksumIEEE([]byte(key))) % this.bucketNum
}

func (this *ShardingStorage) findStorage(key string) IStorage {
	return this.buckets[this.mapToIndex(key)]
}

func (this *ShardingStorage) Set(key string, value []byte, ttl time.Duration) error {
	return this.findStorage(key).Set(key, value, ttl)
}

func (this *ShardingStorage) Get(key string) (*KVItem, error) {
	return this.findStorage(key).Get(key)
}

func (this *ShardingStorage) Delete(key string) error {
	return this.findStorage(key).Delete(key)
}
