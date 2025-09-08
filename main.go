package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/HarrisonZz/web_server_in_go/internal/cache"
	"github.com/HarrisonZz/web_server_in_go/internal/deps"
	"strconv"
	"github.com/HarrisonZz/web_server_in_go/internal/server"
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {

	port := getenv("PORT", "8080")

	addr := getenv("REDIS_ADDR", "127.0.0.1:6379")
	pwd  := getenv("REDIS_PASSWORD", "")
	dbS  := getenv("REDIS_DB", "0")
	db, _ := strconv.Atoi(dbS)

	cache := cache.NewRedisCache(addr, pwd, db)
	defer cache.Close()
	// 讀取埠號（預設 8080）
	r := server.NewRouter(deps.Deps{Cache: cache})

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
