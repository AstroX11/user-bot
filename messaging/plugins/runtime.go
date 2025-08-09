package plugins

import (
	"fmt"
	"time"

	"bot/messaging"
	"bot/messaging/helpers"
	"bot/types"
	"bot/utils"

	"go.mau.fi/whatsmeow/types/events"
)

func init() {
	messaging.RegisterCommand(&types.Command{
		Name:     "runtime",
		Category: "System",
		FromMe:   false,
		IsGroup:  false,
		Handler:  Runtime,
	})
}

func Runtime(msg *events.Message, _ []string) {
	uptime := helpers.FormatRuntime(time.Since(helpers.StartedAt))
	response := fmt.Sprintf("```\nRuntime: %s\n```", uptime)
	_, _ = utils.SendMessage(msg.Info.Chat, response)
}
