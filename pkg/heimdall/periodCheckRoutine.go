package heimdall

import (
	"context"
	"github.com/davidbanham/human_duration"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/jedib0t/go-pretty/v6/table"
	"slices"
	"sort"
	"strings"
	"time"
)

func PeriodicCheckRoutine(ticker time.Ticker, cli *client.Client, ctx context.Context) {
	for {
		<-ticker.C

		containerList, err := cli.ContainerList(ctx, types.ContainerListOptions{
			All: AllContainers,
		})
		if err != nil {
			panic(err)
		}

		t := table.NewWriter()
		t.AppendHeader(table.Row{"#", "Container", "Status"})
		t.SetStyle(table.StyleRounded)

		sort.Slice(containerList, func(i, j int) bool {
			return strings.Compare(containerList[i].State, containerList[j].State) < 0
		})

		slices.Reverse(containerList)

		for i, listContainer := range containerList {
			inspectedContainer, err := cli.ContainerInspect(ctx, listContainer.ID)
			if err != nil {
				panic(err)
			}

			var health = inspectedContainer.State.Status

			if inspectedContainer.State.Health != nil {
				health = inspectedContainer.State.Health.Status
			}

			startTimeStr := inspectedContainer.State.StartedAt

			// Parse the start time string into a time.Time object
			startTime, err := time.Parse(time.RFC3339, startTimeStr)
			if err != nil {
				Fatal(err.Error())
			}

			// Calculate the uptime
			uptime := time.Since(startTime)

			c := Container{
				ID:     listContainer.ID,
				Name:   strings.Split(inspectedContainer.Name, "/")[1],
				Image:  listContainer.Image,
				Uptime: human_duration.String(uptime, human_duration.Second),
				Health: HealthStatus(health),
			}

			t.AppendRow(table.Row{i, c.Name, c.Health})
			t.AppendSeparator()
		}

		Provider.SendPeriodicContainerStatusUpdate(t.Render())
	}
}
