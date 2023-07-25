package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	wav "github.com/go-audio/wav"
)

func Whisper(modelpath, filepath, lang string) (string, error) {
	var result string
	var data []float32

	bChinese := lang == "zh"

	// Load the model
	model, err := whisper.New(modelpath)
	if err != nil {
		return result, err
	}
	defer model.Close()

	fh, err := os.Open(filepath)
	if err != nil {
		return result, err
	}
	defer fh.Close()

	// Decode the WAV file - load the full buffer
	dec := wav.NewDecoder(fh)
	if buf, err := dec.FullPCMBuffer(); err != nil {
		return result, err
	} else if dec.SampleRate != whisper.SampleRate {
		return result, fmt.Errorf("unsupported sample rate: %d", dec.SampleRate)
	} else if dec.NumChans != 1 {
		return result, fmt.Errorf("unsupported number of channels: %d", dec.NumChans)
	} else {
		data = buf.AsFloat32Buffer().Data
	}

	// Process samples
	context, err := model.NewContext()
	if err != nil {
		return result, err
	}

	if bChinese {
		context.SetLanguage("zh")
		m_t2s = LoadDict("TSCharacters.txt")
	} else if lang == "auto" {
		context.SetLanguage("auto")
	}

	if err := context.Process(data, nil, nil); err != nil {
		return result, err
	}

	text := ""
	// Print out the results
	for {
		segment, err := context.NextSegment()
		if err != nil {
			break
		}
		if bChinese {
			text = T2S(segment.Text)
		}
		fmt.Printf("[%6s->%6s] %s\n", segment.Start, segment.End, text)
		result += text
	}
	return result, nil
}

func main() {
	var modelpath, filepath, lang string
	flag.StringVar(&modelpath, "m", "ggml-small.bin", "model path")
	flag.StringVar(&filepath, "f", "zh.wav", "input WAV file path")
	flag.StringVar(&lang, "l", "zh", "spoken language ('auto' for auto-detect)")
	flag.Parse()
	Whisper(modelpath, filepath, lang)
}
