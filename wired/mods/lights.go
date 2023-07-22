package mods

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/kercre123/wire-os/wired/vars"
)

type LightsMod struct {
	vars.Modification
}

func NewLightsMod() *LightsMod {
	return &LightsMod{}
}

type LightsMod_AcceptJSON struct {
	Freq int `json:"freq"`
}

func (fc *LightsMod) Name() string {
	return "LightsMod"
}

func (fc *LightsMod) Description() string {
	return "Modifies CPU/RAM frequency for faster operation."
}

func (fc *LightsMod) DefaultJSON() any {
	return LightsMod_AcceptJSON{
		// default is balanced
		Freq: 1,
	}
}

func (fc *LightsMod) Save(where string, in string) error {
	fcin, ok := fc.DefaultJSON().(LightsMod_AcceptJSON)
	if !ok {
		return errors.New("internal mod error: Save(), DefaultJSON not correct type")
	}
	json.Unmarshal([]byte(in), &fcin)
	saveJson, err := json.Marshal(fcin)
	if err != nil {
		return err
	}
	os.MkdirAll(vars.GetModDir(fc, where), 0777)
	os.WriteFile(vars.GetModDir(fc, where)+"/saved.json", saveJson, 0777)
	return nil
}

func (fc *LightsMod) Load() error {
	fcin, ok := fc.DefaultJSON().(LightsMod_AcceptJSON)
	if !ok {
		return errors.New("internal mod error: Load(), DefaultJSON not correct type")
	}
	file, err := os.ReadFile(vars.GetModDir(fc, "/") + "/saved.json")
	if err != nil {
		defaultJson, _ := json.Marshal(fcin)
		fc.Do("/", string(defaultJson))
		return nil
	}
	json.Unmarshal(file, &fcin)
	doJson, _ := json.Marshal(fcin)
	fc.Do("/", string(doJson))
	return nil
}

func (fc *LightsMod) Accepts() string {
	str, ok := fc.DefaultJSON().(LightsMod_AcceptJSON)
	if !ok {
		log.Fatal("LightsMod Accepts(): not correct type")
	}
	marshedJson, err := json.Marshal(str)
	if err != nil {
		log.Fatal(err)
	}
	return string(marshedJson)
}

func (fc *LightsMod) Do(where string, in string) error {
	fcin, ok := fc.DefaultJSON().(LightsMod_AcceptJSON)
	if !ok {
		return errors.New("internal mod error: Do(), DefaultJSON not correct type")
	}
	err := json.Unmarshal([]byte(in), &fcin)
	if err != nil {
		return err
	}
	fmt.Println(fcin.Freq)
	freq := fcin.Freq
	if freq < 0 || freq > 2 {
		return errors.New("freq must be between 0 and 2")
	}
	var cpufreq string
	var ramfreq string
	var gov string
	switch {
	case freq == 0:
		cpufreq = "533333"
		ramfreq = "400000"
		gov = "interactive"
	case freq == 1:
		cpufreq = "733333"
		ramfreq = "600000"
		gov = "ondemand"
	case freq == 2:
		cpufreq = "1267200"
		ramfreq = "800000"
		gov = "performance"
	}
	fmt.Println(cpufreq + " " + ramfreq + " " + gov)
	fc.Save(where, in)
	return nil
}
