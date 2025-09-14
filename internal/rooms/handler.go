package rooms

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// List Rooms（列出房間 list rooms）
// @Summary 列出所有房間（list all rooms）
// @Tags rooms
// @Produce json
// @Success 200 {array} map[string]any
// @Router /rooms [get]
func List(db *sql.DB) gin.HandlerFunc {
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

// GetRoomsDeviceStates 查詢房間內裝置狀態
// @Summary 查詢房間內裝置狀態
// @Tags rooms
// @Produce json
// @Param room_name query string false "房間名稱 (room name)"
// @Param ids query string false "房間 IDs (逗號分隔)"
// @Success 200 {array} map[string]any
// @Router /rooms/devices/state [get]
func GetDeviceStates(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomName := c.Query("room_name")
		idsStr := c.Query("ids")

		var rows *sql.Rows
		var err error

		if idsStr != "" {
			// 多個房間 id
			idList := strings.Split(idsStr, ",")
			query := `select * from public.get_rooms_current_states(null, $1::int[])`
			rows, err = db.Query(query, "{"+strings.Join(idList, ",")+"}")
		} else if roomName != "" {
			// 指定房間名稱
			query := `select * from public.get_rooms_current_states($1, null)`
			rows, err = db.Query(query, roomName)
		} else {
			// 全部房間
			query := `select * from public.get_rooms_current_states()`
			rows, err = db.Query(query)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close() // 等目前function結束,確保關閉（ensure close）

		var result []map[string]any
		for rows.Next() {
			var (
				roomID     int
				roomName   string
				deviceID   int
				deviceName string
				isOn       sql.NullBool
				lastEvent  sql.NullTime
			)
			if err := rows.Scan(&roomID, &roomName, &deviceID, &deviceName, &isOn, &lastEvent); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			result = append(result, gin.H{
				"room_id":       roomID,
				"room_name":     roomName,
				"device_id":     deviceID,
				"device_name":   deviceName,
				"is_on":         isOn.Bool,
				"last_event_ts": lastEvent.Time,
			})
		}

		c.JSON(http.StatusOK, result)
	}
}

// ChangeDeviceState（改變裝置狀態 change device state）
// @Summary 改變裝置狀態（change device state）
// @Description 新增一筆 on/off 事件，並同步裝置當前狀態（create on/off event and sync device state）
// @Tags devices
// @Accept json
// @Produce json
// @Param id path int true "裝置 id（device id）"
// @Param payload body rooms.ChangeDeviceStateReq true "請求內容（payload）"
// @Success 200 "no content"
// @Failure 400 {object} map[string]any
// @Failure 409 {object} map[string]any "控制被拒（例如總電閘關閉 control gate off）"
// @Failure 500 {object} map[string]any
// @Router /devices/{id}/state [post]
func ChangeDeviceState(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		deviceID, err := strconv.Atoi(idStr)
		if err != nil || deviceID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
			return
		}

		var req ChangeDeviceStateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 呼叫儲存程序（call procedure）
		_, execErr := db.Exec(`call public.apply_device_event($1, $2, $3)`,
			deviceID, req.IsOn, req.Note)

		if execErr != nil {
			c.JSON(http.StatusConflict, gin.H{"error": execErr.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
