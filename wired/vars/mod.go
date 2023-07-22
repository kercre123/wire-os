package vars

import (
	"errors"
	"strings"
)

const (
	SysRoot = "/"
	// for applying mods to updates
	New_SysRoot = "/mnt/"
)

type Modification interface {
	Name() string
	Description() string
	Accepts() string
	DefaultJSON() any
	Save(string, string) error
	// note: Load() runs at init of program
	Load() error
	Do(string, string) error
}

var EnabledMods []Modification

func GetModDir(mod Modification, where string) string {
	//return where + "etc/wired/mods/" + mod.Name()
	return "./modtest/" + mod.Name()
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
		mod.Load()
	}
}
