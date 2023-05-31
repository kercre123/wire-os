package patcher

import (
	"github.com/kercre123/wire-os/pkg/ota-patcher/patches"
	"github.com/kercre123/wire-os/pkg/vars"
)

type OTAPatch struct {
	Name        string
	Description string
	Patch       func(vars.Version, int) error
}

var WireOSPatches []OTAPatch = []OTAPatch{
	{
		Name:        "AddVersion",
		Description: "Puts the desired OTA version and the current time into the build prop, /etc, and the /anki folder.",
		Patch:       patches.AddVersion,
	},
	{
		Name:        "AddCorrectKernelModules",
		Description: "Copies in matching kernel modules for target kernel, so Wi-Fi will work.",
		Patch:       patches.AddCorrectKernelModules,
	},
	{
		Name:        "ProdServerEnv",
		Description: "Switches server env to prod.",
		Patch:       patches.ProdServerEnv,
	},
}
