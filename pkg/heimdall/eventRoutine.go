package heimdall

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"slices"
	"strconv"
	"strings"
)

func EventRoutine(cli *client.Client, ctx context.Context, eventChannel chan ContainerEvent, logChannel chan string, errorChannel chan error) {
	events, errs := cli.Events(ctx, types.EventsOptions{})

	for {
		select {
		case msg, ok := <-events:
			if !ok {
				return
			}
			if msg.Type != "container" || strings.Contains(msg.Action, "exec") {
				continue
			}
			logChannel <- fmt.Sprintf("%+v\n", msg)

			eventsToCapture := []string{"start", "die", "pause", "unpause"}
			if slices.Contains(eventsToCapture, msg.Action) == false && strings.HasPrefix(msg.Action, "health_status") == false {
				continue
			}

			var uptime = "Unknown"
			var health = Unknown
			var exitCode = 0
			var containerError = ""
			var eventType = EventType(msg.Action)
			switch eventType {
			case StartedEvent:
				health = Running
			case DiedEvent:
				health = Exited
				inspectedContainer, err := cli.ContainerInspect(ctx, msg.Actor.ID)
				if msg.Actor.Attributes["exitCode"] != "0" {
					health = Errored
					convertedErrorCode, conversionError := strconv.Atoi(msg.Actor.Attributes["exitCode"])
					if conversionError != nil {
						panic(conversionError)
					}
					exitCode = convertedErrorCode
					containerError = inspectedContainer.State.Error
				}
				if err != nil {
					panic(err)
				}
				uptime = CalculateUptime(inspectedContainer.State.StartedAt, inspectedContainer.State.FinishedAt)
			case PausedEvent:
				health = Paused
				inspectedContainer, err := cli.ContainerInspect(ctx, msg.Actor.ID)
				if err != nil {
					panic(err)
				}
				uptime = CalculateUptime(inspectedContainer.State.StartedAt, inspectedContainer.State.FinishedAt)
			case UnpausedEvent:
				health = Running
			default:
				if strings.HasPrefix(msg.Action, "health_status") {
					health = HealthStatus(strings.Split(msg.Action, ": ")[1])
					eventType = HealthStatusEvent
				}
			}

			event := ContainerEvent{
				Container: Container{
					ID:     msg.Actor.ID,
					Name:   msg.Actor.Attributes["name"],
					Image:  msg.Actor.Attributes["image"],
					Health: health,
					Uptime: uptime,
				},
				Event:    eventType,
				ExitCode: exitCode,
				Error:    containerError,
			}
			eventChannel <- event
		case err, ok := <-errs:
			if !ok {
				return
			}
			errorChannel <- err
		}
	}
}
