package mods

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/kercre123/wire-os/wired/vars"
)

type FreqChange struct {
	vars.Modification
}

func NewFreqChange() *FreqChange {
	return &FreqChange{}
}

var FreqChange_Current FreqChange_AcceptJSON

type FreqChange_AcceptJSON struct {
	Freq int `json:"freq"`
}

func (modu *FreqChange) Name() string {
	return "FreqChange"
}

func (modu *FreqChange) Description() string {
	return "Modifies CPU/RAM frequency for faster operation."
}

func (modu *FreqChange) RestartRequired() bool {
	return false
}

func (modu *FreqChange) DefaultJSON() any {
	return FreqChange_AcceptJSON{
		// default is balanced
		Freq: 1,
	}
}

func (modu *FreqChange) ToFS(to string) {
	// nothing
}

func (modu *FreqChange) Save(where string, in string) error {
	moduin, ok := modu.DefaultJSON().(FreqChange_AcceptJSON)
	if !ok {
		return errors.New("internal mod error: Save(), DefaultJSON not correct type")
	}
	json.Unmarshal([]byte(in), &moduin)
	saveJson, err := json.Marshal(moduin)
	if err != nil {
		return err
	}
	os.MkdirAll(vars.GetModDir(modu, where), 0777)
	os.WriteFile(vars.GetModDir(modu, where)+"saved.json", saveJson, 0777)
	return nil
}

func (modu *FreqChange) Load() error {
	moduin, ok := modu.DefaultJSON().(FreqChange_AcceptJSON)
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
	FreqChange_Current = moduin
	doJson, _ := json.Marshal(moduin)
	modu.Do("/", string(doJson))
	return nil
}

func (modu *FreqChange) Accepts() string {
	str, ok := modu.DefaultJSON().(FreqChange_AcceptJSON)
	if !ok {
		log.Fatal("FreqChange Accepts(): not correct type")
	}
	marshedJson, err := json.Marshal(str)
	if err != nil {
		log.Fatal(err)
	}
	return string(marshedJson)
}

func (modu *FreqChange) Current() string {
	marshalled, _ := json.Marshal(FreqChange_Current)
	return string(marshalled)
}

func (modu *FreqChange) Do(where string, in string) error {
	moduin, ok := modu.DefaultJSON().(FreqChange_AcceptJSON)
	if !ok {
		return errors.New("internal mod error: Do(), DefaultJSON not correct type")
	}
	err := json.Unmarshal([]byte(in), &moduin)
	if err != nil {
		return err
	}
	freq := moduin.Freq
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
	fmt.Println("FreqChange done!: " + cpufreq + " " + ramfreq + " " + gov)
	RunCmd("echo " + cpufreq + " > /sys/devices/system/cpu/cpu0/cpufreq/scaling_max_freq")
	RunCmd("echo disabled > /sys/kernel/debug/msm_otg/bus_voting")
	RunCmd("echo 0 > /sys/kernel/debug/msm-bus-dbg/shell-client/update_request")
	RunCmd("echo 1 > /sys/kernel/debug/msm-bus-dbg/shell-client/mas")
	RunCmd("echo 512 > /sys/kernel/debug/msm-bus-dbg/shell-client/slv")
	RunCmd("echo 0 > /sys/kernel/debug/msm-bus-dbg/shell-client/ab")
	RunCmd("echo active clk2 0 1 max " + ramfreq + " > /sys/kernel/debug/rpm_send_msg/message")
	RunCmd("echo " + gov + " > /sys/devices/system/cpu/cpu0/cpufreq/scaling_governor")
	RunCmd("echo 1 > /sys/kernel/debug/msm-bus-dbg/shell-client/update_request")
	modu.Save(where, in)
	FreqChange_Current = moduin
	return nil
}

func RunCmd(cmd string) ([]byte, error) {
	return exec.Command("/bin/bash", "-c", cmd).Output()
}
