package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	_ "iot-api/docs" // ← 確認 "iot-api" 等於 go.mod 的 module 名稱

	swaggerfiles "github.com/swaggo/files"     // swagger files（swagger files）
	ginSwagger "github.com/swaggo/gin-swagger" // swagger ui handler（swagger ui handler）
)

// @title Smart Home API (smart home api)
// @version 1.0
// @description 裝置/房間狀態管理 API（devices/rooms state management api）
// @host localhost:8080
// @BasePath /

// pingHandler（健康檢查 health check）
// @Summary Ping 測試（ping test）
// @Description 回傳 pong（return pong）
// @Tags health
// @Success 200 {string} string "pong"
// @Router /ping [get]
func pingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// listRooms（列出房間 list rooms）
// @Summary 列出所有房間（list all rooms）
// @Tags rooms
// @Produce json
// @Success 200 {array} map[string]any
// @Router /rooms [get]
func listRooms(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(`select id, name from public.rooms order by id`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var result []map[string]any
		for rows.Next() {
			var id int
			var name string
			if err := rows.Scan(&id, &name); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			result = append(result, gin.H{"id": id, "name": name})
		}
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func main() {
	// 建議改成環境變數；此處先沿用你的 DSN
	connStr := "postgres://postgres:vul3a03sj%2F6u%2C4@10.21.20.28:5432/iot_evergrain?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()

	r.GET("/ping", pingHandler)

	r.GET("/rooms", listRooms(db))

	// swagger ui（只註冊一次）
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
