package api

type LoginRequest struct {
	Name string `json:"name"`
}
type LoginResponse struct {
	UID string `json:"uid"`
}

type HeartbeatRequest struct{}
type HeartbeatResponse struct {
	ServerTime int64 `json:"server_time"`
}
