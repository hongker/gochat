package dto

type IDRequest struct {
	ID string `json:"id"`
}

type EmptyRequest struct{}

type UserResponse User
type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	Location  string `json:"location"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
}
type LoginRequest struct {
	Name string `json:"name"`
}
type LoginResponse struct {
	UID   string `json:"uid"`
	Token string `json:"token"`
}

type ConnectRequest struct {
	UID   string `json:"uid"`
	Token string `json:"token"`
}
type ConnectResponse struct{}

type UserUpdateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Location string `json:"location"`
}
type UserUpdateResponse struct{}

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
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Last  Message `json:"last"`
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
	SessionID   string `json:"session_id"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	CreatedAt   int64  `json:"created_at"`
	Sender      User   `json:"sender"`
}

type ContactQueryRequest struct{}
type ContactQueryResponse struct {
	Items []User `json:"items"`
}

type ChannelQueryRequest struct{}
type ChannelQueryResponse struct {
	Items []Channel `json:"items"`
}

type Channel struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Owner string `json:"owner"`
}
