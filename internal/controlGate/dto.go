package controlGate

type ControlGateResp struct {
	Enabled   bool   `json:"enabled"`
	UpdatedAt string `json:"updated_at"`
}

type ControlGateUpdateReq struct {
	Enabled *bool `json:"enabled" binding:"required"`
}
