package utils

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

var waClient *whatsmeow.Client

func SetClient(c *whatsmeow.Client) {
	waClient = c
}

func SendMessage(jid types.JID, text string) (string, error) {
	if waClient == nil {
		return "", fmt.Errorf("client not initialized")
	}

	msg := &waE2E.Message{
		Conversation: proto.String(text),
	}

	resp, err := waClient.SendMessage(context.Background(), jid, msg)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func EditMessage(jid types.JID, messageID, newText string) error {
	if waClient == nil {
		return fmt.Errorf("client not initialized")
	}

	edit := &waE2E.Message{
		ProtocolMessage: &waE2E.ProtocolMessage{
			Type: waE2E.ProtocolMessage_Type.Enum(14),
			Key: &waCommon.MessageKey{
				RemoteJID: proto.String(jid.String()),
				FromMe:    proto.Bool(true),
				ID:        proto.String(messageID),
			},
			EditedMessage: &waE2E.Message{
				Conversation: proto.String(newText),
			},
		},
	}

	_, err := waClient.SendMessage(context.Background(), jid, edit)
	return err
}
