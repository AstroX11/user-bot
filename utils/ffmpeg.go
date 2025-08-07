package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// ConvertToOpus converts an audio file to Ogg Opus format using ffmpeg.
// It returns the path to the converted file or panics if ffmpeg is not available.
func ConvertToOpus(inputPath string) (string, error) {
	// Check if ffmpeg is installed
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		panic("ffmpeg not found in system path. Please install it to use audio conversion.")
	}

	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".opus.ogg"

	cmd := exec.Command("ffmpeg", "-y", "-i", inputPath, "-c:a", "libopus", "-b:a", "128k", outputPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg conversion failed: %v - %s", err, string(output))
	}

	return outputPath, nil
}

func ConvertToMP3(inputPath string) (string, error) {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		panic("ffmpeg not found in system path. Please install it to use audio conversion.")
	}

	// Don't overwrite source file
	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".converted.mp3"

	cmd := exec.Command("ffmpeg", "-y", "-i", inputPath, "-c:a", "libmp3lame", "-b:a", "192k", outputPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("[ConvertToMP3] ffmpeg error output:\n", string(output))
		return "", fmt.Errorf("ffmpeg mp3 conversion failed: %v", err)
	}

	return outputPath, nil
}

// GetAudioDuration returns the duration of an audio file in seconds.
func GetAudioDuration(filePath string) (uint32, error) {
	// Ensure ffprobe is available
	if _, err := exec.LookPath("ffprobe"); err != nil {
		panic("ffprobe is not installed or not in PATH")
	}

	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		filePath,
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("failed to get duration via ffprobe: %v", err)
	}

	var result struct {
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}

	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		return 0, fmt.Errorf("failed to parse ffprobe output: %v", err)
	}

	// Convert duration string to float then to uint32
	secondsFloat, err := strconv.ParseFloat(result.Format.Duration, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %v", err)
	}

	return uint32(secondsFloat + 0.5), nil // round to nearest second
}
