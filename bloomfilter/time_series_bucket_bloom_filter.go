package bloomfilter

import (
	"fmt"
	"strings"
	"time"
)

// TimeSeriesBucketBloomFilter 时间序列分桶布隆过滤器
type TimeSeriesBucketBloomFilter struct {
	bf         BloomFilter
	keyPrefix  string
	errorRate  float64
	capacity   int64
	bucketUnit string
	ttl        time.Duration
}

// TimeSeriesRow 时间序列数据行
type TimeSeriesRow struct {
	EventTime *time.Time
	EventId   string
}

// NewTimeSeriesBucketBloomFilter 创建时间序列分桶布隆过滤器
// bucketUnit 支持 "day" 和 "hour" 两种分桶单位
func NewTimeSeriesBucketBloomFilter(bloomFilter BloomFilter, keyPrefix string,
	errorRate float64, capacity int64, bucketUnit string, ttl time.Duration) (*TimeSeriesBucketBloomFilter, error) {
	switch bucketUnit {
	case "day", "hour":
	default:
		return nil, fmt.Errorf("unsupported bucket unit: %s", bucketUnit)
	}
	return &TimeSeriesBucketBloomFilter{
		bf:         bloomFilter,
		keyPrefix:  keyPrefix,
		errorRate:  errorRate,
		capacity:   capacity,
		bucketUnit: bucketUnit,
		ttl:        ttl,
	}, nil
}

// GetBucket 根据时间获取桶名称
func (tsbf *TimeSeriesBucketBloomFilter) GetBucket(eventTime *time.Time) string {
	switch tsbf.bucketUnit {
	case "day":
		return eventTime.Format(strings.ReplaceAll(time.DateOnly, "-", ""))
	case "hour":
		return eventTime.Format("2006010215")
	default:
		panic(fmt.Sprintf("Unsupported bucket unit %s", tsbf.bucketUnit))
	}
}

// AddBatch 批量添加时间序列数据
func (tsbf *TimeSeriesBucketBloomFilter) AddBatch(rows []TimeSeriesRow) error {
	keysByBucket := make(map[string][]any)
	for _, row := range rows {
		bucketName := tsbf.GetBucket(row.EventTime)
		keysByBucket[bucketName] = append(keysByBucket[bucketName], row.EventId)
	}

	for bucketName, rows := range keysByBucket {
		err := tsbf.bf.InitAndAddBatch(fmt.Sprintf("%s:%s", tsbf.keyPrefix, bucketName),
			tsbf.errorRate, tsbf.capacity, tsbf.ttl, rows...)
		if err != nil {
			return err
		}
	}

	return nil
}

// ExistsBatch 批量检查时间序列数据是否存在
func (tsbf *TimeSeriesBucketBloomFilter) ExistsBatch(rows []TimeSeriesRow) ([]bool, error) {
	keysByBucket := make(map[string][]any)
	for _, row := range rows {
		bucketName := tsbf.GetBucket(row.EventTime)
		keysByBucket[bucketName] = append(keysByBucket[bucketName], row.EventId)
	}

	existsByKey := make(map[string]bool, len(rows))

	for bucketName, rows := range keysByBucket {
		exists, err := tsbf.bf.ExistsBatch(fmt.Sprintf("%s:%s", tsbf.keyPrefix, bucketName), rows...)
		if err != nil {
			return nil, err
		}
		for i, rowKey := range rows {
			existsByKey[rowKey.(string)] = exists[i]
		}
	}
	res := make([]bool, len(rows))
	for i, row := range rows {
		res[i] = existsByKey[row.EventId]
	}
	return res, nil
}
