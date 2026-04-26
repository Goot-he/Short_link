package MapInit

import (
	"context"
	"new_url/config"
	"new_url/internal/cache"
	"new_url/internal/global"
	"new_url/internal/routes"
	"new_url/pkg/logger"
	"time"
)

func printMsg(err error) {
	if err != nil {
		panic(err)
	}
}

func Init() {
	printMsg(config.InitGlobalConfig()) // 初始化配置文件
	//printMsg(logger.InitLogger(&global.Log))
	global.Log = logger.InitLogger()    // 初始化日志中间件
	printMsg(global.InitGlobalObject()) // 初始化全局引用实例对象
	{
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()             // 必须调用cancel释放资源
		cache.InitBloomFilter(ctx) // 将缓存中的数据加入到布隆过滤器
	}
	routes.InitRouter() // 初始化路由组
}
