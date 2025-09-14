package main

import (
	"log"

	cfgpkg "iot-api/internal/config" // 避免與內部變數同名，取別名（alias）
	dbpkg "iot-api/internal/db"
	"iot-api/internal/router"
	"iot-api/internal/logger"
)

// @title Smart Home API (smart home api)
// @version 1.0
// @description 裝置/房間狀態管理 API（devices/rooms state management api）
// @host localhost:8080
// @BasePath /
func main() {
	
	logger.Init()
	cfg := cfgpkg.Load()
	

	pool, err := dbpkg.Open(cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	r := router.New(pool)

	// 只需啟動一次；使用設定裡的位址（http addr）
	if err := r.Run(cfg.HTTPAddr); err != nil {
		log.Fatal(err)
	}
}
