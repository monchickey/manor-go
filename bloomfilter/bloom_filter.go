package bloomfilter

import "time"

// BloomFilter 布隆过滤器接口定义
type BloomFilter interface {
	Init(key string, errorRate float64, capacity int64) error
	KeyExists(key string) (bool, error)
	KeyExpire(key string, ttl time.Duration) error
	InitAndAddBatch(key string, errorRate float64, capacity int64, ttl time.Duration, item ...any) error
	Add(key, item string) error
	Exists(key, item string) (bool, error)
	AddBatch(key string, item ...any) error
	ExistsBatch(key string, item ...any) ([]bool, error)
}
