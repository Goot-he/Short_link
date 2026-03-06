package service

import (
	"context"
	"errors"
	"fmt"
	"new_url/internal/global"
	"new_url/internal/model"
	"new_url/pkg/base62"
	errormsg "new_url/pkg/errrormsg"
	"time"
)

type Repositories interface {
	CreateUrl(url model.UrlMap) (*model.CreateUrlResponse, error)
	FindShortUrl(shortUrl string) (*model.FindUrlResponse, error)   // 用于重定向
	FindLongUrl(longUrl string) (*model.GetShortUrlResponse, error) // 用于长链接已经存在的情况
	IsAvailableShortUrl(shortUrl string) bool
	IsAvailableLongUrl(longUrl string) bool
}

type Cache interface {
	CreatUrl(ctx context.Context, req model.CreateUrlRequest) error
	FindShortUrl(ctx context.Context, req model.FindUrlRequest) (*model.FindUrlResponse, error)

	SaveToBloomFilter(shortUrl string)
	FindBloomFilter(shortUrl string) bool
}

type ShortCodeGenerator interface {
	CreateShortCode(snowFlakeID int64) string
}

type UrlService struct {
	repositories       Repositories
	cache              Cache
	shortCodeGenerator ShortCodeGenerator
}

func NewUrlService(repo Repositories, cache Cache) *UrlService {
	return &UrlService{
		repositories:       repo,
		cache:              cache,
		shortCodeGenerator: base62.NewShortCodeGenerator(),
	}
}

func (s *UrlService) CreateUrl(ctx context.Context, req model.CreateUrlRequest) (
	*model.CreateUrlResponse, int) {
	var url model.UrlMap
	var err error

	// 判断是否自定义过期时间
	if req.ExpiredAt.IsZero() {
		fmt.Println("用户未自定义过期时间，已将过期时间自动设置为24小时之后")
		req.ExpiredAt = time.Now().UTC().Add(time.Hour * 24)
	} else {
		fmt.Printf("用户自定义过期时间: %v\n", req.ExpiredAt)
	}

	url.LongUrl = req.LongUrl
	url.ShortUrl = req.CustomUrl
	url.ExpiredAt = req.ExpiredAt

	fmt.Printf("[DEBUG] URL Object Created: LongUrl=%s, ShortUrl=%s, ExpiredAt=%v\n", url.LongUrl, url.ShortUrl, url.ExpiredAt)

	// 校验数据 查重与是否合格
	// 检查传入的长链接是否已经存在
	if IsAvailableLongUrl := s.repositories.IsAvailableLongUrl(req.LongUrl); IsAvailableLongUrl {
		fmt.Println("long url is available")
		getShortUrlResponse, err := s.repositories.FindLongUrl(req.LongUrl)
		if err != nil {
			global.Log.Error(err.Error())
			return nil, errormsg.ERROR
		}
		return &model.CreateUrlResponse{
			ShortUrl:  getShortUrlResponse.ShortUrl,
			ExpiredAt: getShortUrlResponse.ExpiredAt,
		}, errormsg.ERROR_LONGURL_AVAILABLE
	}

	// 通过雪花算法生成唯一ID
	url.Id = global.SF.GenerateID()

	// 下面都是在传入的长链接不存在的基础上执行的
	if req.CustomUrl != "" {
		// 这种情况表示用户传入了自定义的short_url
		url.IsCustom = true
		IsAvailableShortUrl := s.repositories.IsAvailableShortUrl(req.CustomUrl)
		if IsAvailableShortUrl {
			// 这里表示自定义的短链已经存在
			return nil, errormsg.ERROR_SHORTURL_AVAILABLE
		}
	} else {
		// 这种情况是用户没有自定义短码 调用函数生成随机短码
		url.IsCustom = false
		url.ShortUrl, err = s.getShortURL(0, url.Id)
		if err != nil {
			global.Log.Error(err.Error())
			return nil, errormsg.ERROR
		}
		if url.ShortUrl == "" { // 如果执行这里表示短码生成失败
			return nil, errormsg.ERROR_GENERATE_SHORTURL_FAIL
		}
	}

	// 关键修复：无论是否自定义，都确保将最终的 ShortUrl 赋值给 req.CustomUrl，以便 Cache 层使用正确的 Key
	req.CustomUrl = url.ShortUrl

	// 调用repo层 将数据存入到数据库
	createUrlResponse, errRepo := s.repositories.CreateUrl(url)
	if errRepo != nil {
		global.Log.Error(errRepo.Error())
		return nil, errormsg.ERROR
	}

	// 调用cache层 将数据存入到缓存
	errCache := s.cache.CreatUrl(ctx, req)
	if errCache != nil {
		global.Log.Error(errCache.Error())
		return nil, errormsg.ERROR
	}

	// 将数据存入到布隆过滤器
	s.cache.SaveToBloomFilter(url.ShortUrl)

	// 返回结果
	return createUrlResponse, errormsg.SUCCESS
}

func (s *UrlService) FindShortUrl(ctx context.Context, req model.FindUrlRequest) (
	*model.FindUrlResponse, error) {
	// 先在布隆过滤器里面判断是否存在
	IsAvailableInBloomFilter := s.cache.FindBloomFilter(req.ShortUrl)
	fmt.Println("shortUrl available in bloomFilter is ", IsAvailableInBloomFilter)

	// 调用cache层 查找缓存
	if IsAvailableInBloomFilter {
		findUrlResponse, errCache := s.cache.FindShortUrl(ctx, req)
		if errCache != nil {
			global.Log.Error(errCache)
		} else {
			if findUrlResponse.LongUrl != "" { // 这里需要确定确实在缓存中查到了目标长链接
				fmt.Println("成功在缓存中查找到目标链接")
				return findUrlResponse, nil
			}
		}
	}

	// 如果缓存中不存在 调用repo层 查找数据库
	findUrlResponse, errRepo := s.repositories.FindShortUrl(req.ShortUrl)
	if errRepo != nil {
		global.Log.Error(errRepo)
	} else if findUrlResponse.LongUrl == "" { // 如果长链接为空表示没有得到目标数据
		fmt.Println("没有在数据库中查找到目标链接")
		return nil, errors.New("shortUrl Not Found")
	}

	// 如果数据库存在并且在缓存中不存在 将数据同步到缓存
	s.cache.SaveToBloomFilter(req.ShortUrl)

	// 返回结果
	return findUrlResponse, nil
}

func (s *UrlService) getShortURL(n int, snowFlakeID int64) (string, error) {
	if n > 5 {
		return "", errors.New("重试次数过多,生成的短链重复率过高")
	}

	shortURL := s.shortCodeGenerator.CreateShortCode(snowFlakeID)
	if !s.repositories.IsAvailableShortUrl(shortURL) {
		return shortURL, nil
	} else {
		return s.getShortURL(n+1, snowFlakeID)
	}
}
