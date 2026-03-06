package model

import (
	//"sync"
	"time"
)

type UrlMap struct {
	Id        int64     `json:"id" gorm:"primaryKey"`
	LongUrl   string    `json:"long_url"`
	ShortUrl  string    `json:"short_url"`
	IsCustom  bool      `json:"is_custom"`
	ExpiredAt time.Time `json:"expired_at" gorm:"column:expired_at"`
}

type CreateUrlRequest struct {
	LongUrl   string    `json:"long_url"`
	CustomUrl string    `json:"custom_url"`
	ExpiredAt time.Time `json:"expired_at"`
}

type CreateUrlResponse struct {
	ShortUrl  string    `json:"short_url"`
	ExpiredAt time.Time `json:"expired_at"`
}

type FindUrlRequest struct {
	ShortUrl string `json:"short_url"`
}

type FindUrlResponse struct {
	LongUrl   string    `json:"long_url"`
	ExpiredAt time.Time `json:"expired_at"`
}

type GetShortUrlResponse struct {
	ShortUrl  string    `json:"short_url"`
	ExpiredAt time.Time `json:"expired_at"`
}
