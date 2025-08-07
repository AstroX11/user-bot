package events

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/AstroX11/user-bot/messaging"
	"github.com/AstroX11/user-bot/sql"
	"github.com/AstroX11/user-bot/utils"
	"go.mau.fi/whatsmeow/types/events"
)

var commandRegex = regexp.MustCompile(`(?i)^\. *([a-z0-9_]+)`)

func Plugins(msg *events.Message) {
	if msg.Message == nil {
		return
	}

	messageText := utils.ExtractText(msg.Message)
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

	match := commandRegex.FindStringSubmatch(messageText)
	if len(match) < 2 {
		return
	}
	cmdName := strings.ToLower(match[1])
	args := strings.Fields(messageText)[1:]

	cmd := messaging.FindCommand(cmdName)
	if cmd != nil {
		cmd.Handler(msg, args)
		return
	}

	suggestion := messaging.SuggestCommand(cmdName)
	if suggestion != "" {
		utils.SendMessage(msg.Info.Chat, fmt.Sprintf("❌ Command `%s` not found. Did you mean `%s%s`?", cmdName, prefix, suggestion))
	} else {
		utils.SendMessage(msg.Info.Chat, fmt.Sprintf("❌ Command `%s` not found.", cmdName))
	}
}
