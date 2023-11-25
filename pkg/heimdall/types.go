package heimdall

type HealthStatus string

const (
	Healthy   HealthStatus = "healthy"
	Unhealthy HealthStatus = "unhealthy"
	Running   HealthStatus = "running"
	Paused    HealthStatus = "paused"
	Exited    HealthStatus = "exited"
	Errored   HealthStatus = "error"
	Unknown   HealthStatus = "unknown"
)

type EventType string

const (
	StartedEvent      EventType = "start"
	DiedEvent         EventType = "die"
	PausedEvent       EventType = "pause"
	UnpausedEvent     EventType = "unpause"
	HealthStatusEvent EventType = "health_status"
)

type Container struct {
	ID     string
	Name   string
	Image  string
	Uptime string
	Health HealthStatus
}

type ContainerEvent struct {
	Container Container
	Event     EventType
	ExitCode  int
	Error     string
}
