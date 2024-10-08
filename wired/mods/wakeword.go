package mods

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kercre123/wire-os/wired/vars"
)

var WakeWordLocation = "/data/data/com.anki.victor/persistent/customWakeWord/wakeword.pmdl"

type WakeWord struct {
	vars.Modification
}

func NewWakeWord() *WakeWord {
	return &WakeWord{}
}

var WakeWord_Current WakeWord_AcceptJSON

type WakeWord_AcceptJSON struct {
	Default bool `json:"default"`
}

func (modu *WakeWord) Name() string {
	return "WakeWord"
}

func (modu *WakeWord) Description() string {
	return "Train a new wake word."
}

func (modu *WakeWord) RestartRequired() bool {
	return true
}

func (modu *WakeWord) DefaultJSON() any {
	return BootAnim_AcceptJSON{
		Default: true,
	}
}

func WakeWord_HTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/mods/wakeword/PrepareListener" {
		vars.StopVic()
		InitListener()
		vars.HTTPSuccess(w, r)
	} else if r.URL.Path == "/api/mods/wakeword/Listen" {
		err := DoListen()
		if err != nil {
			vars.HTTPError(w, r, err.Error())
		} else {
			vars.HTTPSuccess(w, r)
		}
	} else if r.URL.Path == "/api/mods/wakeword/GenWakeWord" {
		if Recind >= 3 && Recind <= 20 {
			if sendWavFilesToServer("/run/wired/wakeword") != nil {
				vars.HTTPError(w, r, "generation error")
			} else {
				vars.HTTPSuccess(w, r)
			}
			return
		} else {
			vars.HTTPError(w, r, "num not in range")
			return
		}
	} else if r.URL.Path == "/api/mods/wakeword/StopListener" {
		StopListener()
		time.Sleep(time.Second)
		vars.StartVic()
		vars.HTTPSuccess(w, r)
	}
}

// func BootAnim_Show() {
// 	// show anim on screen for 10 seconds
// 	cmd := exec.Command("/bin/bash", "-c", "/anki/bin/vic-bootAnim")
// 	vars.StopVic()
// 	go func() {
// 		cmd.Run()
// 	}()
// 	time.Sleep(time.Second * 15)
// 	cmd.Process.Kill()
// 	time.Sleep(time.Second * 1)
// 	vars.StartVic()
// }

func sendWavFilesToServer(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("could not read directory: %w", err)
	}
	var wavFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".wav" {
			wavFiles = append(wavFiles, filepath.Join(dir, file.Name()))
		}
	}
	if len(wavFiles) < 3 || len(wavFiles) > 20 {
		return fmt.Errorf("you must have between 3 and 20 .wav files")
	}
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	for _, wavFile := range wavFiles {
		file, err := os.Open(wavFile)
		if err != nil {
			return fmt.Errorf("could not open wav file: %w", err)
		}
		defer file.Close()
		part, err := writer.CreateFormFile("wavfiles", filepath.Base(wavFile))
		if err != nil {
			return fmt.Errorf("could not create form file: %w", err)
		}
		if _, err := io.Copy(part, file); err != nil {
			return fmt.Errorf("could not copy wav file to form: %w", err)
		}
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("could not close writer: %w", err)
	}
	resp, err := http.Post("http://pvic.xyz:8080/upload", writer.FormDataContentType(), &requestBody)
	if err != nil {
		return fmt.Errorf("could not send post request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned error: %v", resp.Status)
	}
	os.RemoveAll(WakeWordLocation)
	os.MkdirAll(filepath.Dir(WakeWordLocation), 0777)
	outFile, err := os.Create(WakeWordLocation)
	if err != nil {
		return fmt.Errorf("could not create output file: %w", err)
	}
	defer outFile.Close()
	if _, err := io.Copy(outFile, resp.Body); err != nil {
		return fmt.Errorf("could not copy data to file: %w", err)
	}

	return nil
}

func (modu *WakeWord) Save(where string, in string) error {
	return nil
}

func (modu *WakeWord) Load() error {
	return nil
}

func (modu *WakeWord) Accepts() string {
	str, ok := modu.DefaultJSON().(WakeWord_AcceptJSON)
	if !ok {
		log.Fatal("WakeWord Accepts(): not correct type")
	}
	marshedJson, err := json.Marshal(str)
	if err != nil {
		log.Fatal(err)
	}
	return string(marshedJson)
}

func (modu *WakeWord) Current() string {
	marshalled, _ := json.Marshal(WakeWord_Current)
	return string(marshalled)
}

func (modu *WakeWord) Do(where string, in string) error {
	return nil
}
