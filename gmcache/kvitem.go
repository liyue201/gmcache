package gmcache

import (
	"time"
)

type KVItem struct {
	key    string
	value  []byte
	expire time.Duration
}

func (this *KVItem) expired() bool {
	if int64(this.expire) < time.Now().UnixNano() {
		return true
	}
	return false
}
