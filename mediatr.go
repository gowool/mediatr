package mediatr

func Clear() {
	ClearRequestHandlers()
	ClearPipelineBehaviors()
	ClearNotificationHandlers()
}
