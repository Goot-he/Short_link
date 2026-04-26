package logmiddle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"math"
	"new_url/internal/global"
	"os"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置计算消耗时间
		starTtime := time.Now()
		c.Next() //直接调用下一个中间件

		stopTime := time.Since(starTtime)
		spendTime := fmt.Sprintf("%d ms", int(math.Ceil(float64(stopTime.Nanoseconds()/1000000.0))))

		// 主机的名字
		hostName, err := os.Hostname()
		if err != nil {
			hostName = "unknown"
		}
		statusCode := c.Writer.Status()    //响应状态码
		clientIp := c.ClientIP()           //用户IP
		userAgent := c.Request.UserAgent() //判断设备（手机、电脑）
		dataSize := c.Writer.Size()        //响应数据大小
		if dataSize < 0 {
			dataSize = 0
		}
		method := c.Request.Method   //http请求方法
		path := c.Request.RequestURI //完整的请求路径

		// 添加字段到这个entry里面
		entry := global.Log.WithFields(logrus.Fields{
			"HostName":  hostName,
			"status":    statusCode,
			"SpendTime": spendTime,
			"Ip":        clientIp,
			"Method":    method,
			"Path":      path,
			"DataSize":  dataSize,
			"Agent":     userAgent,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
		if statusCode >= 500 {
			entry.Error()
		} else if statusCode >= 400 {
			entry.Warn()
		} else {
			entry.Info()
		}
	}
}

// log 日志中间件 在程序开始前初始化一个log全局对象 然后在中间件里面直接使用这个全局日志实例对象
