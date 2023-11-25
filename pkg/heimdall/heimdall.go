package heimdall

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/docker/docker/client"
	"io"
	"strconv"
	"time"
)

var DebugMode bool
var Timeout int
var PeriodicNotifications bool
var PeriodicNotificationInterval int
var AllContainers bool
var Provider IProvider

func init() {
	flag.BoolVar(&DebugMode, "debug", false, "Enable debug mode")

	flag.BoolVar(&PeriodicNotifications, "periodic-notification", false, "Enable periodic notifications about the state of containers")
	flag.IntVar(&PeriodicNotificationInterval, "notification-interval", 60, "Interval in minutes between periodic notifications")
	flag.BoolVar(&AllContainers, "all-containers", false, "Enable periodic notifications for all containers, including stopped ones")

	flag.IntVar(&Timeout, "retry", 10, "Retry in seconds when the docker event stream ends")

	provider := flag.String("provider", "discord", "Provider to use for notifications")
	webhookUrl := flag.String("webhook-url", "", "Webhook URL to use for notifications")

	flag.Parse()

	if PeriodicNotifications {
		Info(fmt.Sprintf("Periodic notifications enabled. Interval: %d minutes\n", PeriodicNotificationInterval))
	} else {
		if PeriodicNotificationInterval != 60 {
			Warn("Periodic notifications disabled. Ignoring notification interval.")
		}
	}

	if *provider == "" {
		Fatal("A provider is required")
	}
	if *webhookUrl == "" {
		Fatal("A webhook URL is required")
	}

	switch *provider {
	case "discord":
		Provider = New(*webhookUrl)
	}

	Info("Heimdall is now running")
	Info("Awaiting events...")
}

func Start() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		panic(err)
	}

	eventChannel := make(chan ContainerEvent)
	logChannel := make(chan string)
	errorChannel := make(chan error)

	go EventRoutine(cli, ctx, eventChannel, logChannel, errorChannel)

	if PeriodicNotifications {
		ticker := time.NewTicker(time.Duration(PeriodicNotificationInterval) * time.Minute)

		go PeriodicCheckRoutine(*ticker, cli, ctx)
	}

	for {
		select {
		case event := <-eventChannel:
			if DebugMode {
				eventJSON, err := json.Marshal(event)
				if err != nil {
					Debug(err.Error())
				}
				Debug(string(eventJSON))
			}

			Provider.SendContainerEventNotification(event)
		case err := <-errorChannel:
			if err == io.EOF {
				Warn(fmt.Sprintf("No containers running. Sleeping for %s seconds and trying again.\n", strconv.Itoa(Timeout)))
				time.Sleep(time.Duration(Timeout * int(time.Second)))
				go EventRoutine(cli, ctx, eventChannel, logChannel, errorChannel)
			} else {
				Fatal(err.Error())
			}
		case logLine := <-logChannel:
			Debug(logLine)
		}
	}

}
