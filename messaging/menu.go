package messaging

import (
	"fmt"

	"github.com/AstroX11/user-bot/sql"
	"github.com/AstroX11/user-bot/utils"
	"go.mau.fi/whatsmeow/types/events"
)

func Help(msg *events.Message) {
	prefix, err := sql.GetPrefix()
	if err != nil || prefix == "" {
		prefix = "."
	}

	response := fmt.Sprintf(`╭─── ᴍᴇɴᴜ ───╮
%sᴘɪɴɢ
%sᴀʟɪᴠᴇ
%sʜᴇʟᴘ
╰──────────────╯`, prefix, prefix, prefix)

	_, _ = utils.SendMessage(msg.Info.Chat, response)
}
