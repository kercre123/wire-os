package patches

import (
	"os"
	"os/exec"

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

func AddHtop(version vars.Version, target int) error {
	vars.PatchLogger("Adding htop to /usr/bin")
	err := cp.Copy(PluginsPath+"AddHtop/htop", WorkPath+"usr/bin/htop")
	if err != nil {
		return err
	}
	return nil
}

func AddWired(version vars.Version, target int) error {
	if target == 3 {
		vars.PatchLogger("This is an orange build. Not installing wired.")
		return nil
	}
	vars.PatchLogger("Building wired... this may take a while...")
	_, err := exec.Command("/bin/bash", "-c", "./wired/build.sh").Output()
	vars.PatchLogger("Installing wired...")
	if err != nil {
		return err
	}
	err = cp.Copy("./wired/build/wired", WorkPath+"usr/bin/wired")
	if err != nil {
		return err
	}
	err = cp.Copy("./wired/wired.service", WorkPath+"etc/systemd/system/multi-user.target.wants/wired.service")
	if err != nil {
		return err
	}
	os.Mkdir(WorkPath+"etc/wired", 0777)
	err = cp.Copy("./wired/modfiles", WorkPath+"etc/wired/mods")
	if err != nil {
		return err
	}
	os.Remove(WorkPath + "etc/iptables/iptables.rules")
	err = cp.Copy("./wired/iptables.rules", WorkPath+"etc/iptables/iptables.rules")
	if err != nil {
		return err
	}
	err = cp.Copy("./wired/webroot", WorkPath+"etc/wired/webroot")
	if err != nil {
		return err
	}
	return nil
}
