package heimdall

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"github.com/robfig/cron/v3"
	"io"
	"os"
	"strconv"
	"time"
)

var DebugMode bool

var PeriodicNotifications bool
var PeriodicNotificationInterval int
var PeriodicNotificationSchedule string
var AllContainers bool

var Retry int

var Provider IProvider

var Hostname string

var IsHealthy = false

func init() {
	fs := ff.NewFlagSet("heimdall")

	fs.BoolVar(&DebugMode, 'd', "debug", "Enable debug mode")
	fs.BoolVar(&PeriodicNotifications, 'n', "periodic-notification", "Enable periodic notifications about the state of containers")
	fs.IntVar(&PeriodicNotificationInterval, 'i', "notification-interval", 60, "Interval in minutes between periodic notifications")
	fs.StringVar(&PeriodicNotificationSchedule, 's', "notification-schedule", "@hourly", "Cron schedule for periodic notifications")
	fs.BoolVar(&AllContainers, 'a', "all-containers", "Enable periodic notifications for all containers, including stopped ones")
	fs.IntVar(&Retry, 'r', "retry", 10, "Retry in seconds when the docker event stream ends")
	fs.StringVar(&Hostname, 'h', "hostname", "", "Hostname to use in notifications. Useful when running multiple instances of Heimdall")
	var provider = fs.StringEnum('p', "provider", "Provider to use for notifications", "discord")
	var webhookUrl = fs.String('w', "webhook-url", "", "Webhook URL to use for notifications")

	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("HEIMDALL"))

	if err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(fs))
		Fatal("Something went wrong while parsing flags: " + err.Error())
	}

	if PeriodicNotifications {
		Info(fmt.Sprintf("Periodic notifications enabled. Interval: %d minutes\n", PeriodicNotificationInterval))
		if AllContainers {
			Info("All containers option specified, sending notifications about all containers, including stopped ones")
		}
		if PeriodicNotificationInterval != 60 {
			Warn("The periodic notification interval parameter is deprecated. Use --notification-schedule instead")
			if PeriodicNotificationSchedule != "" {
				Warn("--notification-schedule is specified. This overwrites the periodic notification interval parameter")
			}
		}
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
	IsHealthy = true
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
		c := cron.New()

		var interval = "@every " + strconv.Itoa(PeriodicNotificationInterval) + "m"
		if PeriodicNotificationSchedule != "" {
			interval = PeriodicNotificationSchedule
		}

		id, err := c.AddFunc(interval, func() { PeriodicCheckRoutine(cli, ctx, c) })
		if err != nil {
			Error(err.Error())
		}
		c.Start()
		c.Entry(id).Job.Run()
	}

	go StartHealthCheckServer()

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
				IsHealthy = false
				Warn(fmt.Sprintf("No containers running. Sleeping for %s seconds and trying again.\n", strconv.Itoa(Retry)))
				time.Sleep(time.Duration(Retry * int(time.Second)))
				IsHealthy = true
				go EventRoutine(cli, ctx, eventChannel, logChannel, errorChannel)
			} else {
				Fatal(err.Error())
			}
		case logLine := <-logChannel:
			Debug(logLine)
		}
	}
}
