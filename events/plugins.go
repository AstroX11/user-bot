package events

import (
	"strings"

	"go.mau.fi/whatsmeow/types/events"

	"github.com/AstroX11/user-bot/messaging"
	"github.com/AstroX11/user-bot/sql"
	"github.com/AstroX11/user-bot/utils"
)

func Plugins(msg *events.Message) {
	if msg.Message == nil {
		return
	}

	messageText := utils.ExtractText(msg.Message)
	if messageText == "" {
		return
	}

	if messageText == "" {
		return
	}

	prefix, err := sql.GetPrefix()
	if err != nil || prefix == "" {
		prefix = "."
	}

	messageText = strings.TrimSpace(messageText)
	if !strings.HasPrefix(messageText, prefix) {
		return
	}

	command := strings.ToLower(strings.Fields(messageText)[0])

	switch command {
	case prefix + "ping":
		messaging.Ping(msg)
	case prefix + "alive":
		messaging.Alive(msg)
	case prefix + "help":
		messaging.Help(msg)
	default:
		messaging.Unknown(msg, command)
	}
}
