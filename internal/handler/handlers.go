package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/HarrisonZz/web_server_in_go/internal/deps"
	"github.com/HarrisonZz/web_server_in_go/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
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

	r.response = fmt.Sprintf("pong from %s", c.ClientIP())
	Replyln(c, http.StatusOK, r.response)
}

type healthzRoute struct {
	response string
}

func (r *healthzRoute) Method() string { return http.MethodGet }
func (r *healthzRoute) Path() string   { return "/healthz" }
func (r *healthzRoute) Handle(c *gin.Context) {
	logger.Info("[" + r.Method() + "] Hardware Status of Server")
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil || len(cpuPercent) == 0 {
		cpuPercent = []float64{0}
	}

	memUsage, err := mem.VirtualMemory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get memory usage",
		})
		return
	}

	r.response = fmt.Sprintf("CPU Percentage = %.2f%%\nMemory Usage = %.2f%%",
		cpuPercent[0], memUsage.UsedPercent,
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
		key = "sys:os-release"
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
	r.response, r.err = os.ReadFile("/etc/os-release")
	if r.err != nil {
		c.String(http.StatusInternalServerError, "read error: %v", r.err)
		return
	}

	// 3. 存回 Redis（失敗不阻斷）
	if cache != nil {
		_ = cache.Set(c, key, r.response, ttl)
	}

	c.Header("X-Cache", "MISS")
	c.Data(http.StatusOK, "text/plain", r.response)
}

type routeWrapper struct {
	path    string
	handler gin.HandlerFunc
}

func (w *routeWrapper) Method() string { return http.MethodGet }
func (w *routeWrapper) Path() string   { return w.path }
func (w *routeWrapper) Handle(c *gin.Context) {
	w.handler(c)
}

func RegisterRoute(path string, handler gin.HandlerFunc) {
	routes = append(routes, &routeWrapper{path: path, handler: handler})
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
