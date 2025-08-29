package handler

import (
	"net/http"
	"strings"
	"os"
	"github.com/gin-gonic/gin"
)

func Replyln(c *gin.Context, status int, msg string) {
    // 確保每次輸出都有換行
    if !strings.HasSuffix(msg, "\n") {
        msg += "\n"
    }
    c.String(status, msg)
}

func Ping(c *gin.Context) {
	Replyln(c, http.StatusOK, "pong")
}

func Healthz(c *gin.Context) {
	Replyln(c, http.StatusOK, "ok")
}

func Read(c *gin.Context) {
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