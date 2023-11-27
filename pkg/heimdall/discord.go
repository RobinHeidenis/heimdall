package heimdall

import (
	"fmt"
	"github.com/gtuk/discordwebhook"
)

type DiscordProvider struct {
	WebhookURL string
}

func New(webhookURL string) *DiscordProvider {
	return &DiscordProvider{
		WebhookURL: webhookURL,
	}
}

func (p DiscordProvider) SendPeriodicContainerStatusUpdate(updateTable string) {
	username := "Heimdall"
	avatarUrl := "https://neellik.com/wp-content/uploads/2022/03/heimdall.png"
	authorName := "Docker Container Status"
	authorIconUrl := "https://www.docker.com/wp-content/uploads/2023/04/cropped-Docker-favicon-32x32.png"
	title := "Periodic status update"
	description := fmt.Sprintf("```\n%s```", updateTable)

	author := discordwebhook.Author{
		Name:    &authorName,
		IconUrl: &authorIconUrl,
	}

	embeds := []discordwebhook.Embed{{
		Author:      &author,
		Title:       &title,
		Description: &description,
	}}

	message := discordwebhook.Message{
		Username:  &username,
		AvatarUrl: &avatarUrl,
		Embeds:    &embeds,
	}

	err := discordwebhook.SendMessage(p.WebhookURL, message)
	if err != nil {
		Error("Failed to send periodic status update through Discord webhook: " + err.Error())
	}
}

func (p DiscordProvider) SendContainerEventNotification(event ContainerEvent) {
	message := makeWebhookMessage(event)
	err := discordwebhook.SendMessage(p.WebhookURL, message)
	if err != nil {
		Fatal(err.Error())
	}
}

func makeWebhookMessage(event ContainerEvent) discordwebhook.Message {
	containerStatus := event.Container.Health
	if event.Container.Health == Exited {
		containerStatus = "Stopped"
	}

	description := fmt.Sprintf("`%s`'s status is now `%s`", event.Container.Name, containerStatus)

	var title string
	switch event.Container.Health {
	case Running:
		title = fmt.Sprintf("Container %s has started", event.Container.Name)
		if event.Event == UnpausedEvent {
			title = fmt.Sprintf("Container %s has been unpaused", event.Container.Name)
		}
	case Exited:
		title = fmt.Sprintf("Container %s has stopped", event.Container.Name)
	case Paused:
		title = fmt.Sprintf("Container %s has been paused", event.Container.Name)
	case Errored:
		title = fmt.Sprintf("Container %s has errored", event.Container.Name)
		description = fmt.Sprintf("`%s`'s status is now `%s` with exit code `%d`", event.Container.Name, "Stopped", event.ExitCode)
	default:
		title = fmt.Sprintf("Container %s is %s", event.Container.Name, event.Container.Health)
	}

	idFieldName := "ID:"
	imageFieldName := "Image:"

	fields := []discordwebhook.Field{
		{Name: &idFieldName, Value: &event.Container.ID},
		{Name: &imageFieldName, Value: &event.Container.Image},
	}

	if event.Container.Uptime != "Unknown" {
		uptimeFieldName := "Uptime:"
		uptimeFieldValue := event.Container.Uptime
		if uptimeFieldValue == "" {
			uptimeFieldValue = "< 1 second"
		}
		fields = append(fields, discordwebhook.Field{Name: &uptimeFieldName, Value: &uptimeFieldValue})
	}

	if event.Error != "" {
		errorFieldName := "Errored:"
		fields = append(fields, discordwebhook.Field{Name: &errorFieldName, Value: &event.Error})
	}

	authorName := "Docker Container Status Update"
	iconUrl := "https://www.docker.com/wp-content/uploads/2023/04/cropped-Docker-favicon-32x32.png"

	author := discordwebhook.Author{
		Name:    &authorName,
		IconUrl: &iconUrl,
	}

	footerText := "Sent by Heimdall"

	footer := discordwebhook.Footer{
		Text: &footerText,
	}

	var embedColor = "0"
	switch event.Container.Health {
	case Healthy:
		embedColor = "65280"
	case Unhealthy:
		embedColor = "15105570"
	case Errored:
		embedColor = "16711680"
	case Running:
		embedColor = "1926125"
	case Paused:
		embedColor = "15844367"
	}

	embeds := []discordwebhook.Embed{{
		Title:       &title,
		Description: &description,
		Fields:      &fields,
		Author:      &author,
		Footer:      &footer,
		Color:       &embedColor,
	}}

	username := "Heimdall"
	avatarUrl := "https://neellik.com/wp-content/uploads/2022/03/heimdall.png"

	return discordwebhook.Message{
		Embeds:    &embeds,
		Username:  &username,
		AvatarUrl: &avatarUrl,
	}
}
