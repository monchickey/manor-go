package bloomfilter

import (
	"context"
	_ "embed"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed bf_init_and_add_batch.lua
var bfInitAndAddBatchScript string

// RedisBloomFilter Redis 布隆过滤器实现
type RedisBloomFilter struct {
	client          *redis.Client
	initBatchScript *redis.Script
}

// NewRedisBloomFilter 创建 Redis 布隆过滤器实例
func NewRedisBloomFilter(client *redis.Client) *RedisBloomFilter {
	return &RedisBloomFilter{
		client:          client,
		initBatchScript: redis.NewScript(bfInitAndAddBatchScript),
	}
}

// Init 初始化布隆过滤器
func (rbf *RedisBloomFilter) Init(key string, errorRate float64, capacity int64) error {
	return rbf.client.BFReserve(context.Background(), key, errorRate, capacity).Err()
}

// KeyExists 检查键是否存在
func (rbf *RedisBloomFilter) KeyExists(key string) (bool, error) {
	exists, err := rbf.client.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}

	if exists == 1 {
		return true, nil
	}
	return false, nil
}

// KeyExpire 设置键过期时间
func (rbf *RedisBloomFilter) KeyExpire(key string, ttl time.Duration) error {
	return rbf.client.Expire(context.Background(), key, ttl).Err()
}

// Add 添加单个元素
func (rbf *RedisBloomFilter) Add(key, item string) error {
	return rbf.client.BFAdd(context.Background(), key, item).Err()
}

// Exists 检查单个元素是否存在
func (rbf *RedisBloomFilter) Exists(key, item string) (bool, error) {
	return rbf.client.BFExists(context.Background(), key, item).Result()
}

// AddBatch 批量添加元素
func (rbf *RedisBloomFilter) AddBatch(key string, item ...any) error {
	if len(item) == 0 {
		return errors.New("the number of elements in the [AddBatch] method cannot be empty")
	}

	return rbf.client.BFMAdd(context.Background(), key, item...).Err()
}

// ExistsBatch 批量检查元素是否存在
func (rbf *RedisBloomFilter) ExistsBatch(key string, item ...any) ([]bool, error) {
	if len(item) == 0 {
		return nil, errors.New("the number of elements in the [ExistsBatch] method cannot be empty")
	}
	return rbf.client.BFMExists(context.Background(), key, item...).Result()
}

// InitAndAddBatch 初始化并批量添加元素（使用 Lua 脚本保证原子性）
func (rbf *RedisBloomFilter) InitAndAddBatch(key string, errorRate float64, capacity int64,
	ttl time.Duration, item ...any) error {
	args := []any{
		strconv.FormatFloat(errorRate, 'f', -1, 64),
		strconv.FormatInt(capacity, 10),
		strconv.FormatInt(int64(ttl.Seconds()), 10),
	}

	args = append(args, item...)

	return rbf.initBatchScript.Eval(context.Background(), rbf.client, []string{key}, args...).Err()
}
