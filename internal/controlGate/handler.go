package controlGate

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetControlGate
// @Summary 取得保全狀態
// @Description 回傳 保全狀態
// @Tags controlGate
// @Success 200 "no content"
// @Router /control-gate [get]
func GetControlGate(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var enabled bool
		var updatedAt time.Time

		err := db.QueryRow(`select enabled, updated_at from control_gate where id=1`).
			Scan(&enabled, &updatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ControlGateResp{
			Enabled:   enabled,
			UpdatedAt: updatedAt.Format(time.RFC3339),
		})
	}
}

// UpdateControlGate 更新保全狀態
// @Summary 更新保全狀態(上班/下班)
// @Description 更新保全狀態
// @Tags controlGate
// @Accept json
// @Produce json
// @Param request body ControlGateUpdateReq true "更新請求"
// @Success 200 {object} ControlGateResp
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /control-gate [post]
func UpdateControlGate(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ControlGateUpdateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec(`update control_gate set enabled=$1, updated_at=now() where id=1`, req.Enabled)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var enabled bool
		var updatedAt time.Time
		_ = db.QueryRow(`select enabled, updated_at from control_gate where id=1`).
			Scan(&enabled, &updatedAt)

		c.JSON(http.StatusOK, ControlGateResp{
			Enabled:   enabled,
			UpdatedAt: updatedAt.Format(time.RFC3339),
		})
	}
}
