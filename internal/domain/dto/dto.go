package dto

type LoginRequest struct {
	Name string `json:"name"`
}
type LoginResponse struct {
	UID   string `json:"uid"`
	Token string `json:"token"`
}

type HeartbeatRequest struct{}
type HeartbeatResponse struct {
	ServerTime int64 `json:"server_time"`
}

type SessionListRequest struct {
}
type SessionListResponse struct {
	Items []Session `json:"items"`
}
type Session struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type MessageSendRequest struct {
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	Target      string `json:"target"`
}
type MessageSendResponse struct{}

type ChannelCreateRequest struct {
	Name string `json:"name"`
}
type ChannelCreateResponse struct {
	ID string
}

type ChannelJoinRequest struct {
	ID string `json:"id"`
}
type ChannelJoinResponse struct{}

type ChannelLeaveRequest struct {
	ID string `json:"id"`
}
type ChannelLeaveResponse struct{}

type ChannelBroadcastRequest struct {
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	Target      string `json:"target"`
}
type ChannelBroadcastResponse struct{}

type MessageQueryRequest struct {
	SessionID string `json:"session_id"`
}
type MessageQueryResponse struct {
	Items []Message `json:"items"`
}
type Message struct {
	ID          string `json:"id"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	CreatedAt   int64  `json:"created_at"`
}
