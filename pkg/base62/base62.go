package base62

import (
	"strings"
)

type ShortCodeGenerator struct {
}

// 定义标准 62 字符集（顺序固定：0-9a-zA-Z）
const characters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewShortCodeGenerator() *ShortCodeGenerator {
	return &ShortCodeGenerator{}
}

// Base62 正确实现：十进制转62进制字符串
func (s *ShortCodeGenerator) CreateShortCode(snowFlakeID int64) string {
	if snowFlakeID == 0 {
		return string(characters[0]) // 处理0的特殊情况
	}
	var result strings.Builder
	for snowFlakeID > 0 {
		remainder := snowFlakeID % 62
		result.WriteByte(characters[remainder]) // 收集低位余数
		snowFlakeID = snowFlakeID / 62
	}
	// 反转结果（将低位在前转为高位在前）
	reversed := reverseString(result.String())
	return reversed
}

// 辅助函数：反转字符串
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
