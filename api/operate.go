package api

const (
	OperateHeartbeat = iota + 1
	OperateConnect
	OperateUpdateUserProfile
	OperateFindUserProfile
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
