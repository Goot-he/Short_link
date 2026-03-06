package repo

import (
	"new_url/internal/global"
	"new_url/internal/model"
	"time"

	"gorm.io/gorm"
)

type UrlRepoUser struct {
	Db *gorm.DB
}

func NewUrlRepoUser() *UrlRepoUser {
	return &UrlRepoUser{
		Db: global.DB,
	}
}

func (r *UrlRepoUser) CreateUrl(url model.UrlMap) (*model.CreateUrlResponse, error) {
	// 将数据存入数据库
	err := r.Db.Create(&url).Error
	if err != nil {
		return nil, err
	} else {
		return &model.CreateUrlResponse{
			ShortUrl:  url.ShortUrl,
			ExpiredAt: url.ExpiredAt,
		}, nil
	}
}

// 传入一个短链接参数 返回一个长链接
func (r *UrlRepoUser) FindShortUrl(shortUrl string) (
	*model.FindUrlResponse, error) {
	// 查找数据库
	var url model.UrlMap
	//  err := r.Db.Select("short_url").Where("short_url = ?", shortUrl).First(&url).Error
	err := r.Db.Where("short_url = ?", shortUrl).First(&url).Error

	if err != nil {
		return nil, err
	} else {
		return &model.FindUrlResponse{
			LongUrl:   url.LongUrl,
			ExpiredAt: url.ExpiredAt,
		}, err
	}
}

// 传入一个长链接参数 根据这个长链接找对应的短链接 用在长链接已经存在的情况下
func (r *UrlRepoUser) FindLongUrl(longUrl string) (
	*model.GetShortUrlResponse, error) {
	var url model.UrlMap
	err := r.Db.Where("long_url = ?", longUrl).First(&url).Error
	if err != nil {
		return nil, err
	} else {
		return &model.GetShortUrlResponse{
			ShortUrl:  url.ShortUrl,
			ExpiredAt: url.ExpiredAt,
		}, nil
	}
}

func (r *UrlRepoUser) IsAvailableShortUrl(shortUrl string) bool {
	var url model.UrlMap
	r.Db.Where("short_url = ?", shortUrl).First(&url)
	if url.Id > 0 && time.Now().Before(url.ExpiredAt) {
		return true //短链存在并且没过期
	} else {
		return false
	}
}

func (r *UrlRepoUser) IsAvailableLongUrl(longUrl string) bool {
	var url model.UrlMap
	r.Db.Where("long_url = ?", longUrl).First(&url)
	if url.Id > 0 && time.Now().Before(url.ExpiredAt) {
		return true //存在并且没过期
	} else {
		return false
	}
}

// DeleteExpiredUrls 删除所有过期的短链接
func (r *UrlRepoUser) DeleteExpiredUrls() error {
	return r.Db.Where("expired_at <= ?", time.Now()).Delete(&model.UrlMap{}).Error
}
