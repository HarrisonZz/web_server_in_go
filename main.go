package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"strings"
	"github.com/gin-gonic/gin"
)

func Replyln(c *gin.Context, status int, msg string) {
    // 確保每次輸出都有換行
    if !strings.HasSuffix(msg, "\n") {
        msg += "\n"
    }
    c.String(status, msg)
}


func main() {
	// 讀取埠號（預設 8080）
	port := getenv("PORT", "8080")

	// 建 Router：預設含 Logger/Recovery 中介層
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.NoRoute(func(c *gin.Context) {
		Replyln(c, http.StatusNotFound, "404 page not found")
	})

	// 簡單路由
	r.GET("/ping", func(c *gin.Context) { Replyln(c, http.StatusOK, "pong") })
	// 健康檢查（給 K8s / LB）
	r.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	r.GET("/read", func(c *gin.Context) {

		data, err := os.ReadFile("/etc/os-release")
		if err != nil {
			c.String(http.StatusInternalServerError, "read error: %v", err)
			return
		}

		c.Data(http.StatusOK, "text/plain; charset=utf-8", data)

	})

	// 服務（含合理超時）
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// 啟動 HTTP（背景）
	go func() {
		log.Printf("listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// 捕捉 Ctrl+C / 容器停止訊號，優雅關機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	}
	log.Println("server exiting")
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
