package middleware

import (
    "iot-api/internal/logger"

    "github.com/gin-gonic/gin"
)

// ErrorLoggerMiddleware 捕捉錯誤並統一紀錄
func ErrorLoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next() // 執行 handler

        if len(c.Errors) > 0 {
            for _, e := range c.Errors {
                logger.Error.Printf("path=%s, error=%s", c.Request.URL.Path, e.Error())
            }
        }
    }
}
