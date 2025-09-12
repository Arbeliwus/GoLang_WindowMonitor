package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ping（健康檢查 health check）
// @Summary Ping 測試（ping test）
// @Description 回傳 pong（return pong）
// @Tags health
// @Success 200 {string} string "pong"
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
