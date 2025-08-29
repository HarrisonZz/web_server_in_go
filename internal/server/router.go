package server

import (
	"github.com/gin-gonic/gin"
	"github.com/HarrisonZz/web_server_in_go/internal/handler"
)

func NewRouter() *gin.Engine {
	// 建 Router：預設含 Logger/Recovery 中介層
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.NoRoute(handler.NoRoute)

	// 簡單路由
	r.GET("/ping", handler.Ping)
	// 健康檢查（給 K8s / LB）
	r.GET("/healthz", handler.Healthz)

	r.GET("/read", handler.Read)

	return r
}