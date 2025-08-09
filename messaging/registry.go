package messaging

import (
	"strings"

	"bot/types"
)

var commandRegistry = make(map[string]*types.Command)

func RegisterCommand(cmd *types.Command) {
	name := strings.ToLower(cmd.Name)
	commandRegistry[name] = cmd
}

func GetAllCommands() map[string]*types.Command {
	return commandRegistry
}

func FindCommand(name string) *types.Command {
	return commandRegistry[strings.ToLower(name)]
}

func SuggestCommand(input string) string {
	bestMatch := ""
	highestScore := 0
	for name := range commandRegistry {
		score := similarity(input, name)
		if score > highestScore {
			highestScore = score
			bestMatch = name
		}
	}
	if highestScore >= 60 {
		return bestMatch
	}
	return ""
}

func similarity(a, b string) int {
	la := strings.ToLower(a)
	lb := strings.ToLower(b)

	same := 0
	for i := 0; i < len(la) && i < len(lb); i++ {
		if la[i] == lb[i] {
			same++
		}
	}
	return (same * 100) / len(lb)
}
