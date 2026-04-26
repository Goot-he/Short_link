package cache

import (
	"context"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/redis/go-redis/v9"
	"new_url/config"
	"new_url/internal/global"
	"new_url/internal/model"
	"time"
)

type Cache struct {
	rdb *redis.Client
	bf  *bloom.BloomFilter
}

func NewCache() *Cache {
	return &Cache{
		rdb: global.Rdb,
		bf:  BF,
	}
}

func (c *Cache) CreatUrl(ctx context.Context, req model.CreateUrlRequest) error {
	// key是short_url value是原始url
	key := req.CustomUrl
	value := req.LongUrl
	remainingDuration := req.ExpiredAt.Sub(time.Now()) // 将time.time类型转换为duration类型  2020 - 01 - 01 --》 准确的时间数字，例如一小时
	err := c.rdb.Set(ctx, key, value, remainingDuration).Err()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (c *Cache) FindShortUrl(ctx context.Context, req model.FindUrlRequest) (*model.FindUrlResponse, error) {
	//var response *model.FindUrlResponse
	response := &model.FindUrlResponse{}
	var err error

	// 1. 获取长链接
	response.LongUrl, err = c.rdb.Get(ctx, req.ShortUrl).Result()
	if err != nil {
		return nil, err
	}

	// 2. 获取剩余过期时间 (TTL)
	ttl, err := c.rdb.TTL(ctx, req.ShortUrl).Result()
	if err != nil {
		// 如果获取 TTL 失败，可以忽略或记录日志，但不应影响返回长链接
		// 这里选择默认返回零值时间
		response.ExpiredAt = time.Time{}
	} else {
		if ttl > 0 {
			response.ExpiredAt = time.Now().Add(ttl)
		} else {
			// -1 (永不过期) 或 -2 (已过期)
			response.ExpiredAt = time.Time{}
		}
	}

	return response, nil
}

func InitBloomFilter(ctx context.Context) {
	// 数据的数量级
	n := config.GlobalCfg.Bloom.N
	// 预初始的错误率
	p := config.GlobalCfg.Bloom.P

	BF = bloom.NewWithEstimates(uint(n), p) //初始化全局的布隆过滤器

	// 先在redis里面预加载已经存在的短链接
	ExitsStrings := LoadShortURL(ctx)
	for _, String := range ExitsStrings {
		BF.Add([]byte(String))
	}
}

func LoadShortURL(ctx context.Context) []string {
	ans := make([]string, 0)
	var cursor uint64 = 0
	for {
		keys, newCursor, err := global.Rdb.Scan(ctx, cursor, "*", 500).Result()
		if err != nil {
			panic(err)
		}
		ans = append(ans, keys...)
		cursor = newCursor
		if cursor == 0 {
			break
		}
	}
	return ans
}
