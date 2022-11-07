package types

type Message struct {
	ID          string `json:"id"`
	SessionID   string `json:"session_id"`
	SessionType string `json:"session_type"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	Target      string `json:"target"`
	Sender      string `json:"sender"`
	CreatedAt   int64  `json:"created_at"`
}
