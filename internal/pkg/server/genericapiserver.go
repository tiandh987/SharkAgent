package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type GenericAPIServer struct {
	*gin.Engine
	// gin 运行模式
	mode string
	// gin 中间件
	middlewares []string
	// 是用于服务器关闭的超时时间。 这指定服务器正常关闭返回之前的超时。
	ShutdownTimeout time.Duration

	insecureServer *http.Server
	// http 服务配置
	InsecureServingInfo *InsecureServingInfo

	secureServer *http.Server
	// TLS 服务配置
	SecureServingInfo *SecureServingInfo

	// 是否开启健康检查 API 接口
	healthz bool
	//
	enableMetrics bool
	//
	enableProfiling bool
}
