package mods

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/kercre123/vector-gobot/pkg/vbody"
	"github.com/maxhawkins/go-webrtcvad"
	"github.com/youpy/go-wav"
)

type AudioChunk struct {
	Audio  []int16
	Active bool
}

var (
	vadthing    *webrtcvad.VAD
	chunkBuffer [][]int16
	bufferMutex sync.Mutex
	micDataBuf  []int16

	AudioChunks []AudioChunk
	Recind      int
)

func InitListener() {
	Recind = 1
	os.RemoveAll("/run/wired/wakeword")
	os.MkdirAll("/run/wired/wakeword", 0777)
	var err error
	vadthing, err = webrtcvad.New()
	if err != nil {
		panic(err)
	}
	vadthing.SetMode(3)
	vbody.ReadOnly = false
	err = vbody.InitSpine()
	if err != nil {
		panic(err)
	}
	vbody.SetLEDs(vbody.LED_OFF, vbody.LED_OFF, vbody.LED_OFF)
}

func DoListen() error {
	AudioChunks = []AudioChunk{}
	kill := make(chan bool)
	var dieListen bool
	go frameGetter(kill)
	vbody.SetLEDs(vbody.LED_BLUE, vbody.LED_BLUE, vbody.LED_BLUE)
	var timeout int
	go func() {
		for {
			time.Sleep(time.Second)
			timeout++
			if timeout >= 10 {
				dieListen = true
				break
			}
		}
	}()
	fmt.Println("Starting?")
	for {
		if dieListen {
			return errors.New("timeout")
		}
		chunk := getNextChunkFromBuffer()
		if chunk == nil {
			fmt.Println("chunk is nil :(")
			continue
		}

		iVoled := increaseVolume(chunk, 17)
		var bufBytes []byte
		binchunk := bytes.NewBuffer(bufBytes)
		binary.Write(binchunk, binary.LittleEndian, iVoled)

		if IsDoneSpeaking(binchunk.Bytes(), iVoled) {
			kill <- true
			indInt := strconv.Itoa(Recind)
			JustDumpAudio(AudioChunks, "/run/wired/wakeword/record"+indInt+".wav")
			Recind++
			vbody.SetLEDs(vbody.LED_OFF, vbody.LED_OFF, vbody.LED_OFF)
			break
		}
	}
	return nil
}

func JustDumpAudio(cunks []AudioChunk, filepath string) {
	var audBuf []int16
	for _, chunk := range cunks {
		audBuf = append(audBuf, chunk.Audio...)
	}
	WriteWAV(audBuf, filepath)
}

func StopListener() {
	fmt.Println("Stopping spine")
	Recind = 0
	AudioChunks = []AudioChunk{}
	vbody.StopSpine()
}

func WriteWAV(audioData []int16, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()
	writer := wav.NewWriter(file, uint32(len(audioData)), 1, 16000, 16)
	samples := make([]wav.Sample, len(audioData))
	for i, sample := range audioData {
		samples[i].Values[0] = int(sample)
	}
	if err := writer.WriteSamples(samples); err != nil {
		return fmt.Errorf("failed to write samples: %v", err)
	}

	return nil
}

func frameGetter(kill chan bool) {
	var die bool
	go func() {
		for range kill {
			die = true
			break
		}
	}()
	frameChan := vbody.GetFrameChan()
	for frame := range frameChan {
		if die {
			fmt.Println("frame getter is dying")
			break
		}
		smashed := smashPCM(frame.MicData)
		fullBuf, _, isFull := fillBuf(smashed)
		if isFull {
			fillBuffer(fullBuf)
		}
	}
}

func fillBuffer(data []int16) {
	bufferMutex.Lock()
	defer bufferMutex.Unlock()

	chunkBuffer = append(chunkBuffer, data)
}

func fillBuf(in []int16) (full []int16, leftover []int16, filled bool) {
	for i, inny := range in {
		micDataBuf = append(micDataBuf, inny)
		if len(micDataBuf) == 320 {
			mbuf := micDataBuf
			micDataBuf = []int16{}
			return mbuf, in[i+1:], true
		}
	}
	return nil, micDataBuf, false
}

// Retrieve the next chunk from the buffer, waiting until it's not nil
func getNextChunkFromBuffer() []int16 {
	for {
		bufferMutex.Lock()
		if len(chunkBuffer) > 0 {
			chunk := chunkBuffer[0]
			chunkBuffer = chunkBuffer[1:]
			bufferMutex.Unlock()
			return chunk
		}
		bufferMutex.Unlock()
		time.Sleep(5 * time.Millisecond) // wait a bit before checking again to prevent busy-waiting
	}
}

var activeCount int
var inactiveCount int

func IsDoneSpeaking(chunk320 []byte, chunkInt []int16) bool {
	// technically lower than 16000 but whatevs
	active, err := vadthing.Process(16000, chunk320)
	if err != nil {
		panic(err)
	}
	AudioChunks = append(AudioChunks, AudioChunk{
		Audio:  chunkInt,
		Active: active,
	})
	fmt.Println(active)
	if active {
		inactiveCount = 0
		activeCount++
	} else {
		inactiveCount++
		if inactiveCount == 15 {
			if activeCount >= 15 {
				activeCount = 0
				inactiveCount = 0
				return true
			} else {
				activeCount = 0
			}
		}
	}
	return false
}

func increaseVolume(input []int16, factor int16) []int16 {
	output := make([]int16, len(input))

	for i, sample := range input {
		newSample := int32(sample) * int32(factor)
		if newSample > math.MaxInt16 {
			newSample = math.MaxInt16
		} else if newSample < math.MinInt16 {
			newSample = math.MinInt16
		}

		output[i] = int16(newSample)
	}

	return output
}

func smashPCM(input []int16) []int16 {
	if len(input) != 320 {
		panic("gotta be 320 m8")
	}

	output := make([]int16, 80)

	// Extract only the 2nd channel
	for i := 0; i < 80; i++ {
		output[i] = input[i*4]
	}

	return output
}
