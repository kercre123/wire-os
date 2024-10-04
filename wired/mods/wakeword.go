package mods

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/kercre123/wire-os/wired/vars"
)

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
		vars.HTTPError(w, r, "not impl yet")
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
