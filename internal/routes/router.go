package routes

import (
	"new_url/config"
	"new_url/docs"
	"new_url/internal/api"
	"new_url/internal/cache"
	"new_url/internal/global"
	"new_url/internal/repo"
	"new_url/internal/service"
	"new_url/pkg/middlewares/cor"
	"new_url/pkg/middlewares/logmiddle"
	"time"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "new_url/docs" // 必须导入生成的 docs 包（替换为实际路径）

	"github.com/gin-gonic/gin"
)

// @title 短链接服务 API 文档
// @version 1.0
// @description 基于 Gin 实现的短链接生成/查询服务
// @host localhost:8082
// @BasePath
func InitRouter() {
	gin.SetMode(config.GlobalCfg.Server.Mode)
	r := gin.New() // 获取一个gin引擎实例

	docs.SwaggerInfo.Title = "短链接服务 API"
	docs.SwaggerInfo.Description = "基于 Gin 的短链接生成/查询服务"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8082"
	docs.SwaggerInfo.BasePath = "/"

	r.Use(gin.Recovery())
	r.Use(cor.Cors())
	r.Use(logmiddle.Logger())

	// 静态资源托管
	r.Static("/assets", "./frontend/dist/assets")
	r.StaticFile("/", "./frontend/dist/index.html")

	repository := repo.NewUrlRepoUser()
	// periodic cleanup worker: delete expired DB records
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if err := repository.DeleteExpiredUrls(); err != nil {
				global.Log.Error(err)
			}
		}
	}()

	cacheRdb := cache.NewCache()
	MapService := service.NewUrlService(repository, cacheRdb)
	UrlHandler := api.NewUrlHandler(MapService)

	// 修改为通配符路由，匹配如 http://localhost:8082/Abc12
	// 这样可以直接从路径参数中获取短码
	r.GET("/:code", UrlHandler.FindShortUrl)

	apiV1 := r.Group("/api/v1")
	{
		apiV1.POST("/url/create", UrlHandler.CreateUrl)
	}
	// 注册 Swagger 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(config.GlobalCfg.Server.HttpPort); err != nil {
		panic(err)
	}
}
