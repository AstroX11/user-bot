package utils

import (
	"context"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"

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

func SendImage(jid types.JID, imagePath string, caption string) (string, error) {
	if waClient == nil {
		return "", fmt.Errorf("client not initialized")
	}

	data, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %v", err)
	}

	uploaded, err := waClient.Upload(context.Background(), data, whatsmeow.MediaImage)
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %v", err)
	}

	fileInfo, err := os.Stat(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %v", err)
	}

	mimeType := mime.TypeByExtension(filepath.Ext(imagePath))
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	msg := &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(fileInfo.Size())),
			Caption:       proto.String(caption),
		},
	}

	resp, err := waClient.SendMessage(context.Background(), jid, msg)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func SendVideo(jid types.JID, videoPath string, caption string) (string, error) {
	if waClient == nil {
		return "", fmt.Errorf("client not initialized")
	}

	data, err := os.ReadFile(videoPath)
	if err != nil {
		return "", fmt.Errorf("failed to read video file: %v", err)
	}

	uploaded, err := waClient.Upload(context.Background(), data, whatsmeow.MediaVideo)
	if err != nil {
		return "", fmt.Errorf("failed to upload video: %v", err)
	}

	fileInfo, err := os.Stat(videoPath)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %v", err)
	}

	mimeType := mime.TypeByExtension(filepath.Ext(videoPath))
	if mimeType == "" {
		mimeType = "video/mp4"
	}

	msg := &waE2E.Message{
		VideoMessage: &waE2E.VideoMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(fileInfo.Size())),
			Caption:       proto.String(caption),
			Seconds:       proto.Uint32(0),
		},
	}

	resp, err := waClient.SendMessage(context.Background(), jid, msg)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func SendDocument(jid types.JID, documentPath string, fileName string) (string, error) {
	if waClient == nil {
		return "", fmt.Errorf("client not initialized")
	}

	data, err := os.ReadFile(documentPath)
	if err != nil {
		return "", fmt.Errorf("failed to read document file: %v", err)
	}

	uploaded, err := waClient.Upload(context.Background(), data, whatsmeow.MediaDocument)
	if err != nil {
		return "", fmt.Errorf("failed to upload document: %v", err)
	}

	fileInfo, err := os.Stat(documentPath)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %v", err)
	}

	if fileName == "" {
		fileName = filepath.Base(documentPath)
	}

	mimeType := mime.TypeByExtension(filepath.Ext(documentPath))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	msg := &waE2E.Message{
		DocumentMessage: &waE2E.DocumentMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(fileInfo.Size())),
			FileName:      proto.String(fileName),
		},
	}

	resp, err := waClient.SendMessage(context.Background(), jid, msg)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func SendAudio(jid types.JID, audioPath string, isVoiceNote bool) (string, error) {
	if waClient == nil {
		return "", fmt.Errorf("client not initialized")
	}

	var convertedPath string
	var err error

	if isVoiceNote {
		convertedPath, err = ConvertToOpus(audioPath)
	} else {
		convertedPath, err = ConvertToMP3(audioPath)
	}
	if err != nil {
		return "", fmt.Errorf("audio conversion failed: %v", err)
	}
	defer os.Remove(convertedPath)

	data, err := os.ReadFile(convertedPath)
	if err != nil {
		return "", fmt.Errorf("failed to read converted audio file: %v", err)
	}

	uploaded, err := waClient.Upload(context.Background(), data, whatsmeow.MediaAudio)
	if err != nil {
		return "", fmt.Errorf("failed to upload audio: %v", err)
	}

	fileInfo, err := os.Stat(convertedPath)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %v", err)
	}

	duration, err := GetAudioDuration(convertedPath)
	if err != nil {
		return "", fmt.Errorf("failed to get audio duration: %v", err)
	}

	var mimeType string
	var waveform []byte

	if isVoiceNote {
		mimeType = "audio/ogg; codecs=opus"

		if pcm, err := ReadWaveFile(audioPath); err == nil {
			waveform = GenerateWaveform(pcm, 192)
		}
	} else {
		mimeType = "audio/mpeg"
	}

	msg := &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(fileInfo.Size())),
			PTT:           proto.Bool(isVoiceNote),
			Seconds:       proto.Uint32(duration),
			Waveform:      waveform,
		},
	}

	resp, err := waClient.SendMessage(context.Background(), jid, msg)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func SendSticker(jid types.JID, stickerPath string) (string, error) {
	if waClient == nil {
		return "", fmt.Errorf("client not initialized")
	}

	data, err := os.ReadFile(stickerPath)
	if err != nil {
		return "", fmt.Errorf("failed to read sticker file: %v", err)
	}

	uploaded, err := waClient.Upload(context.Background(), data, whatsmeow.MediaImage)
	if err != nil {
		return "", fmt.Errorf("failed to upload sticker: %v", err)
	}

	fileInfo, err := os.Stat(stickerPath)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %v", err)
	}

	mimeType := mime.TypeByExtension(filepath.Ext(stickerPath))
	if mimeType == "" {
		mimeType = "image/webp"
	}

	isAnimated := strings.Contains(strings.ToLower(filepath.Ext(stickerPath)), "gif") ||
		strings.Contains(strings.ToLower(mimeType), "gif")

	msg := &waE2E.Message{
		StickerMessage: &waE2E.StickerMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(fileInfo.Size())),
			IsAnimated:    proto.Bool(isAnimated),
		},
	}

	resp, err := waClient.SendMessage(context.Background(), jid, msg)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func SendMediaFromBytes(jid types.JID, data []byte, mediaType whatsmeow.MediaType, fileName string, caption string) (string, error) {
	if waClient == nil {
		return "", fmt.Errorf("client not initialized")
	}

	uploaded, err := waClient.Upload(context.Background(), data, mediaType)
	if err != nil {
		return "", fmt.Errorf("failed to upload media: %v", err)
	}

	mimeType := getMimeTypeFromFileName(fileName)
	var msg *waE2E.Message

	switch mediaType {
	case whatsmeow.MediaImage:
		msg = &waE2E.Message{
			ImageMessage: &waE2E.ImageMessage{
				URL:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String(mimeType),
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(data))),
				Caption:       proto.String(caption),
			},
		}

	case whatsmeow.MediaVideo:
		msg = &waE2E.Message{
			VideoMessage: &waE2E.VideoMessage{
				URL:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String(mimeType),
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(data))),
				Caption:       proto.String(caption),
				Seconds:       proto.Uint32(0),
			},
		}

	case whatsmeow.MediaDocument:
		msg = &waE2E.Message{
			DocumentMessage: &waE2E.DocumentMessage{
				URL:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String(mimeType),
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(data))),
				FileName:      proto.String(fileName),
			},
		}

	case whatsmeow.MediaAudio:
		msg = &waE2E.Message{
			AudioMessage: &waE2E.AudioMessage{
				URL:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String(mimeType),
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(data))),
				Seconds:       proto.Uint32(0),
				PTT:           proto.Bool(strings.Contains(strings.ToLower(fileName), "voice")),
			},
		}

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
func getMimeTypeFromFileName(fileName string) string {
	ext := filepath.Ext(fileName)
	mimeType := mime.TypeByExtension(ext)
	if mimeType != "" {
		return mimeType
	}

	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".mp4":
		return "video/mp4"
	case ".avi":
		return "video/avi"
	case ".mov":
		return "video/quicktime"
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".mp3":
		return "audio/mpeg"
	case ".ogg":
		return "audio/ogg"
	case ".wav":
		return "audio/wav"
	default:
		return "application/octet-stream"
	}
}

func ExtractText(msg *waE2E.Message) string {
	switch {
	case msg.GetConversation() != "":
		return msg.GetConversation()
	case msg.ExtendedTextMessage != nil:
		return msg.ExtendedTextMessage.GetText()
	case msg.ImageMessage != nil && msg.ImageMessage.Caption != nil:
		return msg.ImageMessage.GetCaption()
	case msg.VideoMessage != nil && msg.VideoMessage.Caption != nil:
		return msg.VideoMessage.GetCaption()
	case msg.ProtocolMessage != nil && msg.ProtocolMessage.EditedMessage != nil:
		return ExtractText(msg.ProtocolMessage.EditedMessage)
	}
	return ""
}

func GetQuotedMessage(msg *waE2E.ContextInfo) {

}
