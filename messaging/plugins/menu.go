package plugins

import (
	"fmt"
	"runtime"
	"sort"
	"time"

	"bot/config"
	"bot/messaging"
	"bot/messaging/helpers"
	"bot/sql"
	"bot/types"
	"bot/utils"

	"go.mau.fi/whatsmeow/types/events"
)

func init() {
	messaging.RegisterCommand(&types.Command{
		Name:     "menu",
		Category: "System",
		FromMe:   false,
		IsGroup:  false,
		Handler:  Help,
	})
}

func Help(msg *events.Message, _ []string) {
	prefix, err := sql.GetPrefix()
	if err != nil || prefix == "" {
		prefix = "."
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
	uptime := time.Since(helpers.StartedAt)
	day := time.Now().Weekday().String()
	date := time.Now().Format("01/02/2006")
	tm := time.Now().Format("15:04:05")

	allCommands := messaging.GetAllCommands()
	categorized := make(map[string][]string)

	for _, cmd := range allCommands {
		categorized[cmd.Category] = append(categorized[cmd.Category], prefix+utils.FancyText(cmd.Name))
	}

	var categories []string
	for cat := range categorized {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	infoBlock := fmt.Sprintf("```╭─── %s ────\n", botName) +
		fmt.Sprintf("│ User: %s\n", pushName) +
		fmt.Sprintf("│ Owner: %s\n", owner) +
		fmt.Sprintf("│ Plugins: %d\n", len(allCommands)) +
		fmt.Sprintf("│ Mode: %s\n", mode) +
		fmt.Sprintf("│ Uptime: %s\n", helpers.FormatRuntime(uptime)) +
		fmt.Sprintf("│ Platform: %s\n", runtime.GOOS) +
		fmt.Sprintf("│ Ram: %s\n", helpers.FormatMemUsage()) +
		fmt.Sprintf("│ Day: %s\n", day) +
		fmt.Sprintf("│ Date: %s\n", date) +
		fmt.Sprintf("│ Time: %s\n", tm) +
		fmt.Sprintf("│ Go: %s\n", runtime.Version()) +
		"╰─────────────```\n"

	cmdBlock := ""
	for _, cat := range categories {
		cmds := categorized[cat]
		fancyCat := utils.FancyText(cat)

		top := fmt.Sprintf("╭─── %s ───", fancyCat)

		lines := make([]string, len(cmds))
		for i, cmd := range cmds {
			lines[i] = fmt.Sprintf("│ %d %s", i+1, cmd)
		}

		cmdBlock += "```" + top + "\n"
		for _, line := range lines {
			cmdBlock += line + "\n"
		}
		cmdBlock += "╰───────```\n"
	}

	utils.SendMessage(msg.Info.Chat, infoBlock+cmdBlock)
}
