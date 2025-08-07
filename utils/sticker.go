package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
)

func ImageToWebp(inputPath, outputPath string) error {
	args := []string{
		"-i", inputPath,
		"-vcodec", "libwebp",
		"-vf", "scale='min(320,iw)':min'(320,ih)':force_original_aspect_ratio=decrease,fps=15,pad=320:320:-1:-1:color=white@0.0",
		"-lossless", "1",
		"-preset", "default",
		"-an",
		"-y", outputPath,
	}
	return runFFmpeg(args)
}

func VideoToWebp(inputPath, outputPath string, durationSec int) error {
	args := []string{
		"-i", inputPath,
		"-vcodec", "libwebp",
		"-vf", "scale='min(320,iw)':min'(320,ih)':force_original_aspect_ratio=decrease,fps=15,pad=320:320:-1:-1:color=white@0.0",
		"-loop", "0",
		"-ss", "00:00:00",
		fmt.Sprintf("-t=%02d", durationSec),
		"-preset", "default",
		"-an",
		"-y", outputPath,
	}
	return runFFmpeg(args)
}

func runFFmpeg(args []string) error {
	cmd := exec.Command("ffmpeg", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %v: %s", err, stderr.String())
	}
	return nil
}

func InjectEXIFMeta(inputPath, outputPath, packName, author string) error {
	json := fmt.Sprintf(`{"sticker-pack-id":"com.astrox11.user-bot","sticker-pack-name":"%s","sticker-pack-publisher":"%s","emojis":["⚡"]}`, packName, author)

	jsonBytes := []byte(json)
	exifPayload := append([]byte("II*\x00\x08\x00\x00\x00"), jsonBytes...)

	var buf bytes.Buffer
	buf.WriteString("EXIF")
	if err := binary.Write(&buf, binary.BigEndian, uint32(len(exifPayload))); err != nil {
		return err
	}
	buf.Write(exifPayload)

	exifChunk := buf.Bytes()

	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	var output []byte
	inserted := false

	for i := 0; i < len(inputData); {
		if i+8 > len(inputData) {
			break
		}

		chunkType := string(inputData[i : i+4])
		chunkSize := binary.BigEndian.Uint32(inputData[i+4 : i+8])
		chunkEnd := i + 8 + int(chunkSize)

		output = append(output, inputData[i:chunkEnd]...)

		if chunkType == "VP8X" && !inserted {
			output = append(output, exifChunk...)
			inserted = true
		}

		i = chunkEnd
	}

	if !inserted {
		return fmt.Errorf("VP8X chunk not found in file")
	}

	return os.WriteFile(outputPath, output, 0644)
}
