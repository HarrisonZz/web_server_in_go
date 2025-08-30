package server

import (
	"github.com/HarrisonZz/web_server_in_go/internal/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	// 建 Router：預設含 Logger/Recovery 中介層
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.NoRoute(handler.NoRoute)

	routeMap := handler.GetRoutes()
	for path, h := range routeMap {
		r.GET(path, h)
	}

	return r
}
