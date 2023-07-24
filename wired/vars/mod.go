package vars

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	SysRoot = "/"
	// for applying mods to updates
	Update_SysRoot  = "/mnt/"
	VectorResources = "anki/data/assets/cozmo_resources/"
)

type Modification interface {
	Name() string
	Description() string
	Accepts() string
	DefaultJSON() any
	Save(string, string) error
	// note: Load() runs at init of program
	Load() error
	// current settings of mod
	Current() string
	// fs root
	ToFS(string)
	RestartRequired() bool
	Do(string, string) error
}

type BaseModification struct {
	Modification
	ModName            string
	ModDescription     string
	VicRestartRequired bool
}

func (bc *BaseModification) Name() string {
	return bc.ModName
}

func (bc *BaseModification) Description() string {
	return bc.ModDescription
}

func (bc *BaseModification) RestartRequired() bool {
	return bc.VicRestartRequired
}

var EnabledMods []Modification

func GetModDir(mod Modification, where string) string {
	return where + "etc/wired/mods/" + mod.Name() + "/"
	//return "./modtest/" + mod.Name() + "/"
}

func FindMod(name string) (Modification, error) {
	for index, mod := range EnabledMods {
		if strings.TrimSpace(name) == mod.Name() {
			return EnabledMods[index], nil
		}
	}
	return nil, errors.New("mod not found")
}

func InitMods() {
	for _, mod := range EnabledMods {
		fmt.Println("Loading " + mod.Name() + "...")
		mod.Load()
	}
}

func RestartVic() {
	exec.Command("/bin/bash", "-c", "systemctl stop anki-robot.target").Output()
	time.Sleep(time.Second * 4)
	exec.Command("/bin/bash", "-c", "systemctl start anki-robot.target").Output()
	time.Sleep(time.Second * 3)
}
