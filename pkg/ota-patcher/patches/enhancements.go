package patches

import (
	"os"

	"github.com/kercre123/wire-os/pkg/vars"
	cp "github.com/otiai10/copy"
)

const (
	BootAnimLocation = PluginsPath + "AddCustomBootAnim/boot_anim.raw"
	AnkiInitLocation = PluginsPath + "UpCPUFreq/ankiinit"
)

func AddCustomBootAnim(version vars.Version, target int) error {
	vars.PatchLogger("Replacing boot animation")
	os.Remove(WorkPath + "anki/data/assets/cozmo_resources/config/engine/animations/boot_anim.raw")
	err := cp.Copy(BootAnimLocation, WorkPath+"anki/data/assets/cozmo_resources/config/engine/animations/boot_anim.raw")
	if err != nil {
		return err
	}
	return nil
}

func UpCPUFreq(version vars.Version, target int) error {
	vars.PatchLogger("Replacing ankiinit script")
	os.Remove(WorkPath + "etc/initscripts/ankiinit")
	err := cp.Copy(AnkiInitLocation, WorkPath+"etc/initscripts/ankiinit")
	if err != nil {
		return err
	}
	return nil
}
