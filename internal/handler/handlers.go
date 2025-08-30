package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

var getRoutes = map[string]gin.HandlerFunc{
	"/ping":    ping,
	"/healthz": healthz,
	"/os":      osInfo,
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

func ping(c *gin.Context) {
	clientIP := c.ClientIP()
	msg := fmt.Sprintf("pong from %s", clientIP)
	Replyln(c, http.StatusOK, msg)
}

func healthz(c *gin.Context) {
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

func osInfo(c *gin.Context) {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		Replyln(c, http.StatusInternalServerError, "read error: "+err.Error())
		return
	}
	if len(data) == 0 || data[len(data)-1] != '\n' {
		data = append(data, '\n')
	}
	c.Data(http.StatusOK, "text/plain; charset=utf-8", data)
}

func NoRoute(c *gin.Context) {
	Replyln(c, http.StatusNotFound, "404 page not found")
}
