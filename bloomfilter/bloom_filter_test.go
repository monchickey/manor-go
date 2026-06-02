package bloomfilter_test

import (
	"testing"
	"time"

	"github.com/monchickey/manor-go/v2/bloomfilter"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestBloomFilter(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Username: "redis",
		Password: "redis-123456",
		DB:       0,
	})
	defer rdb.Close()

	bf := bloomfilter.NewRedisBloomFilter(rdb)
	err := bf.Init("bf-test", 0.001, 10000)
	assert.Nil(t, err)
	e, err := bf.KeyExists("bf-test")
	assert.Nil(t, err)
	assert.True(t, e)
	e1, err := bf.KeyExists("bf-test1")
	assert.Nil(t, err)
	assert.False(t, e1)

	err = bf.Add("bf-test", "hhh")
	assert.Nil(t, err)
	e, err = bf.Exists("bf-test", "hhh")
	assert.Nil(t, err)
	assert.True(t, e)

	err = bf.AddBatch("bf-test", "hhh1", "hhh2", "hhh3")
	assert.Nil(t, err)
	es, err := bf.ExistsBatch("bf-test", "hhh1", "hhh2", "hhh3", "hhh4", "hhh5")
	assert.Nil(t, err)
	assert.True(t, es[0])
	assert.True(t, es[1])
	assert.True(t, es[2])
	assert.False(t, es[3])
	assert.False(t, es[4])

	err = bf.KeyExpire("bf-test", time.Second*3)
	assert.Nil(t, err)
	time.Sleep(5 * time.Second)
	e, err = bf.KeyExists("bf-test")
	assert.Nil(t, err)
	assert.False(t, e)
}
