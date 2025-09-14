package router

import (
	"database/sql"
	"net/http"

	_ "iot-api/docs" // swagger docs（swagger docs）

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"     // swagger files（swagger files）
	ginSwagger "github.com/swaggo/gin-swagger" // swagger ui handler（swagger ui handler）

	"iot-api/internal/health"
	"iot-api/internal/rooms"
)

func New(db *sql.DB) *gin.Engine {
	r := gin.Default()

	// health（健康檢查 health check）
	r.GET("/ping", health.Ping)

	// rooms（房間 rooms）
	r.GET("/rooms", rooms.List(db))

	// 房間 + 裝置狀態
	r.GET("/rooms/devices/state", rooms.GetDeviceStates(db))


	// swagger ui（swagger ui）
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 可加一個根路由提示（optional）
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "smart-home api running. see /swagger")
	})

	return r
}
