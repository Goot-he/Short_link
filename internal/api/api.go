package api

import (
	"context"
	"net/http"
	"new_url/internal/global"
	"new_url/internal/model"
	errormsg "new_url/pkg/errrormsg"

	"github.com/gin-gonic/gin"
)

type UrlService interface {
	CreateUrl(ctx context.Context, req model.CreateUrlRequest) (*model.CreateUrlResponse, int)
	FindShortUrl(ctx context.Context, req model.FindUrlRequest) (*model.FindUrlResponse, error)
}
type UrlHandler struct {
	urlService UrlService
}

func NewUrlHandler(service UrlService) *UrlHandler {
	return &UrlHandler{
		urlService: service,
	}
}

// @Summary 创建短链接
// @Description 接收长链接、可选的自定义短码与过期时间，返回短码与过期时间
// @Tags URL
// @Accept json
// @Produce json
// @Param body body model.CreateUrlRequest true "请求体"
// @Success 200 {object} model.CreateUrlResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/url/create [post]
func (h *UrlHandler) CreateUrl(c *gin.Context) {

	// 获取传入的参数
	var url model.CreateUrlRequest
	if err := c.ShouldBindJSON(&url); err != nil {
		global.Log.Error(err)
	}

	// 调用server层函数 借助接口调用函数
	rep, errorMsg := h.urlService.CreateUrl(c.Request.Context(), url)

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"status":   errorMsg,
		"errorMsg": errormsg.GetErrorMsg(errorMsg),
		"rep":      rep,
	})
}

// @Summary 短链重定向
// @Description 根据短码重定向到原始长链接
// @Tags URL
// @Param code path string true "短码"
// @Success 301
// @Failure 404 {object} map[string]string
// @Router /{code} [get]
func (h *UrlHandler) FindShortUrl(c *gin.Context) {
	// 1. 尝试从 URL 路径参数获取 (例如 /:code)
	// 这样可以直接从路径参数中获取短码 然后在浏览器中直接根据短码访问原网页
	code := c.Param("code")
	var req model.FindUrlRequest

	if code != "" {
		req.ShortUrl = code
	} else {
		// 2. 如果路径参数为空，尝试从 JSON Body 获取 (兼容旧接口)
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Invalid request parameters"})
			return
		}
	}

	// 调用service层函数
	rep, err := h.urlService.FindShortUrl(c.Request.Context(), req)
	if err != nil {
		global.Log.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"status": "error",
			"msg": "Short URL not found",
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, rep.LongUrl)
}
