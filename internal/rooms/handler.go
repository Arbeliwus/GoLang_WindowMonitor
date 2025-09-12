package rooms

import (
	"database/sql"
	"net/http"

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
