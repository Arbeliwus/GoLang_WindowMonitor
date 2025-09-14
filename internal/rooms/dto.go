package rooms

// ChangeStateReq 改變裝置狀態的請求格式
type ChangeDeviceStateReq struct {
    IsOn bool    `json:"is_on" binding:"required"`
    Note *string `json:"note"`
}