package logger

import (
	"log"
	"os"
	"path/filepath"
)

var Logger *log.Logger

// Init 初始化 logger
func Init(logPath string) error {
	// 確保目錄存在
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 打開檔案
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// 建立 logger
	Logger = log.New(f, "APP_LOG: ", log.LstdFlags|log.Llongfile)
	return nil
}

// Info 寫入 info log
func Info(msg string) {
	if Logger != nil {
		Logger.Println("[INFO]", msg)
	}
}

// Error 寫入 error log
func Error(msg string) {
	if Logger != nil {
		Logger.Println("[ERROR]", msg)
	}
}
