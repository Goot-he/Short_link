package global

import (
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"new_url/database/Db"
	//"new_url/internal/cache"
	"new_url/pkg/snowflake"
)

var (
	machineID = int64(1) // 机械ID
)

var (
	err error
	DB  *gorm.DB
	Rdb *redis.Client
	SF  *snowflake.SnowflakeIDGenerator
	Log = logrus.New()
)

// 初始化全局对象
func InitGlobalObject() error {
	DB = Db.NewDB()
	Rdb = Db.NewRedis()
	SF, err = snowflake.NewSnowflakeIDGenerator(machineID)
	return err
}
