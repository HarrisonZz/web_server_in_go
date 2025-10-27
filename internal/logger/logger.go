package logger

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"
)

var Logger *log.Logger

type LogEntry struct {
	Timestamp string `json:"@timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Host      string `json:"path,omitempty"`
}

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
	Logger = log.New(f, "", 0)
	return nil
}

func input(msg string, level string) {

	if Logger != nil {

		host, _ := os.Hostname()
		entry := LogEntry{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Level:     level,
			Message:   msg,
			Host:      host,
		}
		jsonData, err := json.Marshal(entry)
		if err != nil {
			Logger.Printf(`{"level":"error","message":"failed to marshal log: %s"}`, err)
			return
		}
		Logger.Println(string(jsonData))
	}

}

// Info 寫入 info log
func Info(msg string) {

	input(msg, "info")

}

// Error 寫入 error log
func Error(msg string) {

	input(msg, "info")

}

func Warn(msg string) {

	input(msg, "info")

}
