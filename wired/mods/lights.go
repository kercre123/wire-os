package mods

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/kercre123/wire-os/wired/vars"
	cp "github.com/otiai10/copy"
)

type RainbowLights struct {
	vars.Modification
}

func NewRainbowLights() *RainbowLights {
	return &RainbowLights{}
}

var RainbowLights_Current RainbowLights_AcceptJSON

type RainbowLights_AcceptJSON struct {
	Enabled bool `json:"enabled"`
}

func (modu *RainbowLights) Name() string {
	return "RainbowLights"
}

func (modu *RainbowLights) Description() string {
	return "Makes the backpack/cube lights rainbow."
}

func (modu *RainbowLights) RestartRequired() bool {
	return true
}

func (modu *RainbowLights) DefaultJSON() any {
	return RainbowLights_AcceptJSON{
		Enabled: true,
	}
}

func (modu *RainbowLights) Save(where string, in string) error {
	moduin, ok := modu.DefaultJSON().(RainbowLights_AcceptJSON)
	if !ok {
		return errors.New("internal mod error: Save(), DefaultJSON not correct type")
	}
	json.Unmarshal([]byte(in), &moduin)
	saveJson, err := json.Marshal(moduin)
	if err != nil {
		return err
	}
	os.MkdirAll(vars.GetModDir(modu, where), 0777)
	os.WriteFile(vars.GetModDir(modu, where)+"/saved.json", saveJson, 0777)
	return nil
}

func (modu *RainbowLights) Load() error {
	moduin, ok := modu.DefaultJSON().(RainbowLights_AcceptJSON)
	if !ok {
		return errors.New("internal mod error: Load(), DefaultJSON not correct type")
	}
	file, err := os.ReadFile(vars.GetModDir(modu, "/") + "saved.json")
	if err != nil {
		defaultJson, _ := json.Marshal(moduin)
		modu.Do("/", string(defaultJson))
		return nil
	}
	json.Unmarshal(file, &moduin)
	RainbowLights_Current = moduin
	return nil
}

func (modu *RainbowLights) Accepts() string {
	str, ok := modu.DefaultJSON().(RainbowLights_AcceptJSON)
	if !ok {
		log.Fatal("RainbowLights Accepts(): not correct type")
	}
	marshedJson, err := json.Marshal(str)
	if err != nil {
		log.Fatal(err)
	}
	return string(marshedJson)
}

func (modu *RainbowLights) Current() string {
	marshalled, _ := json.Marshal(RainbowLights_Current)
	return string(marshalled)
}

func (modu *RainbowLights) Do(where string, in string) error {
	moduin, ok := modu.DefaultJSON().(RainbowLights_AcceptJSON)
	if !ok {
		return errors.New("internal mod error: Do(), DefaultJSON not correct type")
	}
	err := json.Unmarshal([]byte(in), &moduin)
	if err != nil {
		return err
	}
	fmt.Println("RainbowLights Do(): " + fmt.Sprint(moduin.Enabled))
	var lightsFolder string
	destFolder := where + vars.VectorResources + "config/engine/lights"
	if moduin.Enabled {
		lightsFolder = vars.GetModDir(modu, where) + "rainbow"
	} else {
		lightsFolder = vars.GetModDir(modu, where) + "orig"
	}
	os.RemoveAll(destFolder)
	cp.Copy(lightsFolder, destFolder)
	modu.Save(where, in)
	RainbowLights_Current = moduin
	return nil
}
