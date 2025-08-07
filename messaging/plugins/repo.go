package plugins

import (
	"github.com/AstroX11/user-bot/messaging"
	"github.com/AstroX11/user-bot/types"
	"github.com/AstroX11/user-bot/utils"
	"go.mau.fi/whatsmeow/types/events"
)

func init() {
	messaging.RegisterCommand(&types.Command{
		Name:     "repo",
		FromMe:   false,
		Category: "misc",
		Handler:  repo,
	})
}

func repo(msg *events.Message, _ []string) {
	 utils.SendImage(msg.Info.Sender, "./resources/logo.png", "Simple User WhatsAppBot\nhttps://github.com/AstroX11/user-bot")
}
