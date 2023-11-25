package heimdall

type IProvider interface {
	SendPeriodicContainerStatusUpdate(updateTable string)
	SendContainerEventNotification(event ContainerEvent)
}
