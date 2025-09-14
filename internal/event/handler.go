package event

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ListEvents 查詢事件
// @Summary 查詢事件
// @Description 可以依時間區間和房間過濾
// @Tags events
// @Produce json
// @Param start query string false "開始時間 (RFC3339 格式)"
// @Param end query string false "結束時間 (RFC3339 格式)"
// @Param room_id query int false "房間 ID"
// @Success 200 {array} EventResp
// @Router /events [get]
func ListEvents(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		startStr := c.Query("start")
		endStr := c.Query("end")
		roomIDStr := c.Query("room_id")

		var start, end *time.Time
		var roomID *int

		if startStr != "" {
			if t, err := time.Parse(time.RFC3339, startStr); err == nil {
				start = &t
			}
		}
		if endStr != "" {
			if t, err := time.Parse(time.RFC3339, endStr); err == nil {
				end = &t
			}
		}
		if roomIDStr != "" {
			if id, err := strconv.Atoi(roomIDStr); err == nil {
				roomID = &id
			}
		}

		rows, err := db.Query(`
            select e.id, e.device_id, d.device_name, r.id, r.name, e.is_open, e.event_timestamp
            from public.devices_events e
            join public.devices d on e.device_id = d.id
            join public.rooms r on d.room_id = r.id
            where 
                ($1::timestamptz is null or e.event_timestamp >= $1) and
                ($2::timestamptz is null or e.event_timestamp <= $2) and
                ($3::int is null or r.id = $3)
            order by e.event_timestamp desc
        `, start, end, roomID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var result []EventResp
		for rows.Next() {
			var ev EventResp
			if err := rows.Scan(&ev.EventID, &ev.DeviceID, &ev.DeviceName, &ev.RoomID, &ev.RoomName, &ev.IsOpen, &ev.EventTimestamp); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			result = append(result, ev)
		}

		c.JSON(http.StatusOK, result)
	}
}
