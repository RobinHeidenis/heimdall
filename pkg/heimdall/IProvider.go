package heimdall

import "time"

type IProvider interface {
	SendPeriodicContainerStatusUpdate(updateTable string, nextInvocation time.Time)
	SendContainerEventNotification(event ContainerEvent)
}
