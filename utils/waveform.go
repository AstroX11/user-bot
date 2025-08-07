package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-audio/wav"
)

// ReadWaveFile attempts to read a WAV file directly,
// and falls back to ffmpeg decoding for other formats.
func ReadWaveFile(filePath string) ([]int, error) {
	ext := filepath.Ext(filePath)

	// Prefer native decode for .wav files
	if ext == ".wav" {
		samples, err := readWavDirect(filePath)
		if err == nil {
			return samples, nil
		}
	}

	// Fallback: use ffmpeg to decode audio to 16-bit PCM mono
	return readPcmViaFfmpeg(filePath)
}

// readWavDirect uses the go-audio/wav library for WAV files.
func readWavDirect(filePath string) ([]int, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := wav.NewDecoder(f)
	if !decoder.IsValidFile() {
		return nil, fmt.Errorf("invalid WAV file")
	}

	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, err
	}

	return buf.Data, nil
}

// readPcmViaFfmpeg uses ffmpeg to extract 16-bit mono PCM from any audio file.
func readPcmViaFfmpeg(filePath string) ([]int, error) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		panic("ffmpeg not installed or not in PATH")
	}

	cmd := exec.Command("ffmpeg",
		"-i", filePath,
		"-f", "s16le", // 16-bit signed PCM
		"-acodec", "pcm_s16le", // codec
		"-ac", "1", // mono
		"-ar", "16000", // 16 kHz
		"-hide_banner",
		"-loglevel", "quiet",
		"pipe:1", // output to stdout
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg PCM extraction failed: %v", err)
	}

	// Read as int16 PCM data
	data := out.Bytes()
	sampleCount := len(data) / 2
	samples := make([]int, sampleCount)

	for i := 0; i < sampleCount; i++ {
		sample := int16(binary.LittleEndian.Uint16(data[i*2:]))
		samples[i] = int(sample)
	}

	return samples, nil
}

func GenerateWaveform(samples []int, count int) []byte {
	if len(samples) == 0 {
		return nil
	}
	if len(samples) < count {
		count = len(samples)
	}

	step := len(samples) / count
	waveform := make([]byte, count)

	for i := 0; i < count; i++ {
		sumSquares := 0
		for j := 0; j < step; j++ {
			idx := i*step + j
			if idx >= len(samples) {
				break
			}
			v := samples[idx]
			sumSquares += v * v
		}

		// Calculate RMS amplitude
		rms := int(math.Sqrt(float64(sumSquares) / float64(step)))

		// Normalize to 0–100 scale (based on max int16 = 32767)
		if rms > 32767 {
			rms = 32767
		}
		waveform[i] = byte((rms * 100) / 32767)
	}

	return waveform
}
