package api

const (
	OperateHeartbeat = iota + 1
	OperateLogin
	OperateUpdateProfile
	OperateListSession
	OperateSendMessage

	OperateCreateChannel
	OperateJoinChannel
	OperateLeaveChannel
	OperateBroadcastChannel
	OperateQueryMessage
)

const (
	OperatePushMessage = 101
)
