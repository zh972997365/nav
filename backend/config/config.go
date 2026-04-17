package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	// 获取当前文件的绝对路径
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Failed to get current file path")
	}

	// 构建 .env 文件的绝对路径（当前文件所在目录的上一级）
	envPath := filepath.Join(filepath.Dir(filename), "..", ".env")

	if err := godotenv.Load(envPath); err != nil {
		log.Println("No .env file found at:", envPath)
	}
}
