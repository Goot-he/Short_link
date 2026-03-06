package Db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"new_url/config"
)

func NewDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GlobalCfg.DataBase.DbUser,
		config.GlobalCfg.DataBase.DbPassword,
		config.GlobalCfg.DataBase.DbHost,
		config.GlobalCfg.DataBase.DbPort,
		config.GlobalCfg.DataBase.DbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
