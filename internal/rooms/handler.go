package rooms

import (
	"database/sql"
	"net/http"
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
				roomID   int
				roomName string
				deviceID int
				deviceName string
				isOn sql.NullBool
				lastEvent sql.NullTime
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