package heimdall

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/robfig/cron/v3"
	"slices"
	"sort"
	"strings"
)

func PeriodicCheckRoutine(cli *client.Client, ctx context.Context, cron *cron.Cron) {
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

		if inspectedContainer.State.Health != nil && health != "exited" {
			health = inspectedContainer.State.Health.Status
		}

		t.AppendRow(table.Row{i, strings.Split(inspectedContainer.Name, "/")[1], HealthStatus(health)})
		t.AppendSeparator()
	}

	nextInvocation := cron.Entries()[0].Next

	Provider.SendPeriodicContainerStatusUpdate(t.Render(), nextInvocation)
}
