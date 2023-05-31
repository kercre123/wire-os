package patches

import (
	"os"

	"github.com/kercre123/wire-os/pkg/vars"
	cp "github.com/otiai10/copy"
)

func AddRsync(version vars.Version, target int) error {
	vars.PatchLogger("Adding rsync to /bin")
	os.Remove(WorkPath + "bin/rsync")
	err := cp.Copy(PluginsPath+"AddRsync/rsync", WorkPath+"bin/rsync")
	if err != nil {
		return err
	}
	return nil
}

func CustomUpdateEngine(version vars.Version, target int) error {
	envEntries := `
	UPDATE_ENGINE_ENABLED=True
	UPDATE_ENGINE_ALLOW_DOWNGRADE=True
	UPDATE_ENGINE_BASE_URL=http://wire.my.to:81/` + vars.Targets[target] + `-stable/
	UPDATE_ENGINE_BASE_URL_LATEST=http://wire.my.to:81/` + vars.Targets[target] + `-unstable/
	`
	vars.PatchLogger("Adding custom update-engine to /anki/bin (saving original to /anki/bin/orig-update-engine)")
	cp.Copy(WorkPath+"anki/bin/update-engine", WorkPath+"anki/bin/orig-update-engine")
	err := cp.Copy(PluginsPath+"CustomUpdateEngine/update-engine", WorkPath+"anki/bin/update-engine")
	if err != nil {
		return err
	}
	vars.PatchLogger("Adding custom update-engine env to /anki/etc")
	err = os.WriteFile(WorkPath+"anki/etc/update-engine.env", []byte(envEntries), 0777)
	if err != nil {
		return err
	}
	return nil
}

func AddNano(version vars.Version, target int) error {
	vars.PatchLogger("Adding nano to /usr/bin")
	err := cp.Copy(PluginsPath+"AddNano/nano", WorkPath+"usr/bin/nano")
	if err != nil {
		return err
	}
	vars.PatchLogger("Adding libncursesw and libtinfo to /lib")
	err = cp.Copy(PluginsPath+"AddNano/libncursesw.so.5", WorkPath+"usr/lib/libncursesw.so.5")
	if err != nil {
		return err
	}
	err = cp.Copy(PluginsPath+"AddNano/libtinfo.so.5", WorkPath+"lib/libtinfo.so.5")
	if err != nil {
		return err
	}
	err = cp.Copy(PluginsPath+"AddNano/libncurses.so.5", WorkPath+"lib/libncurses.so.5")
	if err != nil {
		return err
	}
	return nil
}
