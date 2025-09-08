package server

import (
	"github.com/HarrisonZz/web_server_in_go/internal/deps"
	"github.com/HarrisonZz/web_server_in_go/internal/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(d deps.Deps) *gin.Engine {
	// 建 Router：預設含 Logger/Recovery 中介層
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), deps.InjectDeps(d))
	r.NoRoute(handler.NoRoute)

	routeMap := handler.GetRoutes()
	for path, h := range routeMap {
		r.GET(path, h)
	}

	return r
}
