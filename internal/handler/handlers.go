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

var getRoutes = map[string]gin.HandlerFunc{
	"/ping":    ping,
	"/healthz": healthz,
	"/os":      getOSRelease,
}

func GetRoutes() map[string]gin.HandlerFunc {
	// 回傳副本，防止外部修改
	routes := make(map[string]gin.HandlerFunc)
	for key, value := range getRoutes {
		routes[key] = value
	}
	return routes
}

func RegisterRoute(path string, handler gin.HandlerFunc) {
	getRoutes[path] = handler
}

func Replyln(c *gin.Context, status int, msg string) {
	// 確保每次輸出都有換行
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	c.String(status, msg)
}

func getOSRelease(c *gin.Context) {
	logger.Info("[GET] Get OS Release")
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
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		c.String(http.StatusInternalServerError, "read error: %v", err)
		return
	}

	// 3. 存回 Redis（失敗不阻斷）
	if cache != nil {
		_ = cache.Set(c, key, data, ttl)
	}

	c.Header("X-Cache", "MISS")
	c.Data(http.StatusOK, "text/plain", data)
}

func ping(c *gin.Context) {
	logger.Info("[GET] Get Pong Response")
	clientIP := c.ClientIP()
	msg := fmt.Sprintf("pong from %s", clientIP)
	Replyln(c, http.StatusOK, msg)
}

func healthz(c *gin.Context) {
	logger.Info("[GET] Get Node Status")
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

	msg := fmt.Sprintf("CPU Percentage = %.2f%%\nMemory Usage = %.2f%%",
		cpuPercent[0], memUsage.UsedPercent,
	)
	Replyln(c, http.StatusOK, msg)
}

func NoRoute(c *gin.Context) {
	Replyln(c, http.StatusNotFound, "404 page not found")
}
