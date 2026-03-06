package cache

import (
	"github.com/bits-and-blooms/bloom/v3"
	//"new_url/internal/global"
)

var BF *bloom.BloomFilter

// 将数据存入到布隆过滤器
func (c *Cache) SaveToBloomFilter(shortUrl string) {
	c.bf.AddString(shortUrl)
}

// 查找数据是否在布隆过滤器里面
func (c *Cache) FindBloomFilter(shortUrl string) bool {
	IsFind := c.bf.TestString(shortUrl)
	return IsFind
}
