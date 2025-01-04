package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"time"
)

type httpServerInterface interface {
	Timeout(time.Duration) gin.HandlerFunc // 设置http超时时间 	r.Use(middleware.Http().Timeout(1 * time.Second))
}
type httpServer struct{}

//goland:noinspection GoExportedFuncWithUnexportedType
func Http() httpServerInterface {
	return &httpServer{}
}

func (*httpServer) Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建一个带有超时的上下文
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// 使用 goroutine 处理请求，以便不会阻塞主线程
		ch := make(chan struct{})
		go func() {
			c.Next()
			ch <- struct{}{}
		}()

		// 等待请求处理完成或超时
		select {
		case <-ch:
			// 请求处理完成
			return
		case <-ctx.Done():
			// 请求超时
			c.AbortWithStatusJSON(200, gin.H{
				"code":    "8",
				"message": "请求超时",
			})
		}
	}
}
