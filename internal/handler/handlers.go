package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/HarrisonZz/web_server_in_go/internal/deps"
	"github.com/HarrisonZz/web_server_in_go/internal/kubernetes"
	"github.com/HarrisonZz/web_server_in_go/internal/logger"
	"github.com/gin-gonic/gin"
)

var routes = []Route{
	&pingRoute{},
	&healthzRoute{},
	&osInfoRoute{},
}

func GetRoutes() map[string]gin.HandlerFunc {
	// 回傳副本，防止外部修改
	out := make(map[string]gin.HandlerFunc, len(routes))
	for _, r := range routes {

		if r.Method() == http.MethodGet {
			out[r.Path()] = r.Handle
		}
	}
	return out
}

type Route interface {
	Method() string
	Path() string
	Handle(*gin.Context)
}

type pingRoute struct {
	response string
}

func (r *pingRoute) Method() string { return http.MethodGet }
func (r *pingRoute) Path() string   { return "/ping" }
func (r *pingRoute) Handle(c *gin.Context) {
	logger.Info("[" + r.Method() + "] Pong")

	r.response = fmt.Sprintf("pong from %s", kubernetes.NodeInfo.InternalIP)
	Replyln(c, http.StatusOK, r.response)
}

type healthzRoute struct {
	response string
}

func (r *healthzRoute) Method() string { return http.MethodGet }
func (r *healthzRoute) Path() string   { return "/healthz" }
func (r *healthzRoute) Handle(c *gin.Context) {
	logger.Info("[" + r.Method() + "] Hardware Status of Server")
	r.response = fmt.Sprintf(
		"Node: %s\nMemory: %s\nCPU: %s",
		kubernetes.NodeInfo.Name,
		kubernetes.NodeInfo.Memory,
		kubernetes.NodeInfo.CPU,
	)

	Replyln(c, http.StatusOK, r.response)
}

type osInfoRoute struct {
	response []byte
	err      error
}

func (r *osInfoRoute) Method() string { return http.MethodGet }
func (r *osInfoRoute) Path() string   { return "/os" }
func (r *osInfoRoute) Handle(c *gin.Context) {
	logger.Info("[" + r.Method() + "] OS Infomation")
	const (
		key = "sys:os-info"
		ttl = 2 * time.Hour // 資訊幾乎不會變動，可設長一點
	)

	cache := deps.CacheFrom(c)

	// 1. 嘗試讀 Redis
	if cache != nil {
		if b, ok, err := cache.Get(c, key); err == nil && ok {
			c.Header("X-Cache", "HIT")
			c.Data(http.StatusOK, "text/plain", b)
			return
		}
	}

	// 2. MISS → 讀檔
	data := gin.H{
		"node": kubernetes.NodeInfo.Name,
		"os": gin.H{
			"architecture":    kubernetes.NodeInfo.Arch,
			"operatingSystem": kubernetes.NodeInfo.OS,
			"osImage":         kubernetes.NodeInfo.OSImage,
			"kernelVersion":   kubernetes.NodeInfo.Kernel,
		},
	}

	c.Header("X-Cache", "MISS")
	c.IndentedJSON(200, data)

	// 3. 存回 Redis（失敗不阻斷）
	if cache != nil {
		jsonBytes, err := json.Marshal(data)
		if err == nil {
			_ = cache.Set(c, key, jsonBytes, ttl)
		}
	}
}

type routeWrapper struct {
	method  string
	path    string
	handler gin.HandlerFunc
}

func (w *routeWrapper) Method() string { return w.method }
func (w *routeWrapper) Path() string   { return w.path }
func (w *routeWrapper) Handle(c *gin.Context) {
	w.handler(c)
}

func RegisterRoute(method string, path string, handler gin.HandlerFunc) {
	routes = append(routes, &routeWrapper{method: method, path: path, handler: handler})
}

func Replyln(c *gin.Context, status int, msg string) {
	// 確保每次輸出都有換行
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	c.String(status, msg)
}

func NoRoute(c *gin.Context) {
	path := c.Request.URL.Path
	logger.Error(fmt.Sprintf("No route rule for path: %s", path))
	Replyln(c, http.StatusNotFound, "404 page not found")
}
