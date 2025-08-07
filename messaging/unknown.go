package messaging

import (
	"fmt"

	"github.com/AstroX11/user-bot/sql"
	"github.com/AstroX11/user-bot/utils"
	"go.mau.fi/whatsmeow/types/events"
)

func Unknown(msg *events.Message, command string) {
	prefix, err := sql.GetPrefix()
	if err != nil || prefix == "" {
		prefix = "."
	}

	response := fmt.Sprintf("❓ Unknown command: %s\n\nType %shelp to see available commands.", command, prefix)
	_, _ = utils.SendMessage(msg.Info.Chat, response)
}
