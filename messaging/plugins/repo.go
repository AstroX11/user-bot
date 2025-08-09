package plugins

import (
	"bot/messaging"
	"bot/types"
	"bot/utils"

	"go.mau.fi/whatsmeow/types/events"
)

func init() {
	messaging.RegisterCommand(&types.Command{
		Name:     "repo",
		FromMe:   false,
		Category: "misc",
		Handler: func(msg *events.Message, _ []string) {
			utils.SendImage(msg.Info.Chat, "./resources/logo.png", "Simple User WhatsAppBot\nhttps://github.com/AstroX11/user-bot")
		},
	})
}
