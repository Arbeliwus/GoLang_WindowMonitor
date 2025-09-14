package event
import "time"
type EventResp struct {
    EventID       int       `json:"event_id"`
    DeviceID      int       `json:"device_id"`
    DeviceName    string    `json:"device_name"`
    RoomID        int       `json:"room_id"`
    RoomName      string    `json:"room_name"`
    IsOpen        bool      `json:"is_open"`
    EventTimestamp time.Time `json:"event_timestamp"`
}
