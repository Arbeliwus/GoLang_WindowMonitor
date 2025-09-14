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

	r.GET("/ping", health.Ping)
	r.GET("/rooms", rooms.List(db))
	r.GET("/rooms/devices/state", rooms.GetDeviceStates(db))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.POST("/devices/:id/state", rooms.ChangeDeviceState(db))

	// 可加一個根路由提示（optional）
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "smart-home api running. see /swagger")
	})

	return r
}
