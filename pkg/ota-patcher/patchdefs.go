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
		Name:        "AddNano",
		Description: "Adds nano, a text editor, to the filesystem.",
		Patch:       patches.AddNano,
	},
	{
		Name:        "AddWired",
		Description: "Adds wire_d, a daemon for controlling stuff via a webserver.",
		Patch:       patches.AddWired,
	},
	{
		Name:        "AddHtop",
		Description: "Adds htop.",
		Patch:       patches.AddHtop,
	},
	{
		Name:        "AddSysSwitch",
		Description: "Adds a system slot switcher called sysswitch.",
		Patch:       patches.AddSysSwitch,
	},
	{
		Name:        "AddRsync",
		Description: "Adds rsync, which is helpful for data transfer, to the filesystem.",
		Patch:       patches.AddRsync,
	},
	{
		Name:        "UpCPUFreq",
		Description: "Increases the CPU and RAM frequency a bit.",
		Patch:       patches.UpCPUFreq,
	},
	{
		Name:        "PatchMountData",
		Description: "Mounts /data with exec perms",
		Patch:       patches.PatchMountData,
	},
	{
		Name:        "AddSSHKey",
		Description: "Adds regular dev SSH keys (if DDL-era software)",
		Patch:       patches.AddSSHKey,
	},
}
