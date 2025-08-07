package messaging

import (
	"fmt"
	"runtime"
	"time"

	"github.com/AstroX11/user-bot/config"
	"github.com/AstroX11/user-bot/messaging/helpers"
	"github.com/AstroX11/user-bot/sql"
	"github.com/AstroX11/user-bot/utils"
	"go.mau.fi/whatsmeow/types/events"
)

func Help(msg *events.Message) {
	prefix, err := sql.GetPrefix()
	if err != nil || prefix == "" {
		prefix = "."
	}

	commands := []string{"ping", "alive", "help"}

	var cmds []string
	for i, cmd := range commands {
		fancy := utils.FancyText(cmd)
		cmds = append(cmds, fmt.Sprintf("%d. %s%s", i+1, prefix, fancy))
	}

	owner := config.AppConfig.UserName
	if owner == "" {
		owner = "αѕтяσχ11"
	}
	botName := config.AppConfig.BotName
	if botName == "" {
		botName = "ᴜsᴇʀ ʙᴏᴛ"
	}
	pushName := msg.Info.PushName
	if pushName == "" {
		pushName = "Unknown"
	}
	mode := "Public"
	uptime := time.Since(startedAt)
	day := time.Now().Weekday().String()
	date := time.Now().Format("01/02/2006")
	tm := time.Now().Format("15:04:05")

	infoBlock := fmt.Sprintf("```╭─── %s ────\n", botName) +
		fmt.Sprintf("│ User: %s\n", pushName) +
		fmt.Sprintf("│ Owner: %s\n", owner) +
		fmt.Sprintf("│ Plugins: %d\n", len(commands)) +
		fmt.Sprintf("│ Mode: %s\n", mode) +
		fmt.Sprintf("│ Uptime: %s\n", helpers.FormatRuntime(uptime)) +
		fmt.Sprintf("│ Platform: %s\n", runtime.GOOS) +
		fmt.Sprintf("│ Ram: %s\n", helpers.FormatMemUsage()) +
		fmt.Sprintf("│ Day: %s\n", day) +
		fmt.Sprintf("│ Date: %s\n", date) +
		fmt.Sprintf("│ Time: %s\n", tm) +
		fmt.Sprintf("│ Go: %s\n", runtime.Version()) +
		"╰─────────────```\n"

	cmdBlock := "```╭─── Commands ───╮\n"
	for _, item := range cmds {
		cmdBlock += fmt.Sprintf("│ %s\n", item)
	}
	cmdBlock += "╰────────────────╯```"

	utils.SendMessage(msg.Info.Chat, infoBlock+"\n"+cmdBlock)
}
