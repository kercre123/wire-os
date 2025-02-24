package patches

import (
	"bufio"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/kercre123/wire-os/pkg/vars"
	cp "github.com/otiai10/copy"
)

const (
	WorkPath    = "./work/"
	PluginsPath = "./resources/patches/"
)

func AddVersion(version vars.Version, target int) error {
	// revision will eventually be hooked up to github commit
	vars.PatchLogger("Modifying build.prop file")
	// read prop in OTA
	var propLines []string
	var customProps []string
	verAppendage := "dev"
	origProp, err := os.Open(WorkPath + "build.prop")
	if err != nil {
		return err
	}
	propScanner := bufio.NewScanner(origProp)
	propScanner.Split(bufio.ScanLines)
	for propScanner.Scan() {
		text := propScanner.Text()
		if strings.Contains(text, "ro.build.version.release") {
			origProp.Close()
			break
		}
		propLines = append(propLines, text)
	}
	os.Remove(WorkPath + "build.prop")
	currentTime := time.Now()
	formattedTime := currentTime.Format("200601021504")
	vars.PatchLogger("Current time: " + formattedTime)
	vars.PatchLogger("Version: v" + version.Full + verAppendage)
	/*
		ro.revision=anki-ea5d84b_os-689b9cd
		ro.anki.version=2.0.1.6080
		ro.anki.victor.version=2.0.1.6080
		ro.build.fingerprint=v2.0.1.6080-ea5d84b_os2.0.1.6080-689b9cd-202209201956
		ro.build.id=v2.0.1.6080-ea5d84b_os2.0.1.6080-689b9cd-202209201956
		ro.build.display.id=2.0.1.6080
		ro.build.type=production
		ro.build.version.incremental=6080
	*/
	customProps = append(customProps, "ro.build.version.release="+formattedTime)
	customProps = append(customProps, "ro.product.name=Vector")
	customProps = append(customProps, "ro.revision=wire_os")
	customProps = append(customProps, "ro.anki.version="+version.Full+verAppendage)
	customProps = append(customProps, "ro.anki.victor.version="+version.Full)
	customProps = append(customProps, "ro.build.fingerprint="+version.Full+"-wire_os"+version.Full+"-"+verAppendage+"-"+formattedTime)
	customProps = append(customProps, "ro.build.id=v"+version.Full+"-wire_os"+version.Full+"-"+verAppendage+"-"+formattedTime)
	customProps = append(customProps, "ro.build.display.id=v"+version.Full+vars.Targets[target])
	customProps = append(customProps, "ro.build.target="+strconv.Itoa(target))
	customProps = append(customProps, "ro.build.type=development")
	customProps = append(customProps, "ro.build.version.incremental="+version.Increment)
	customProps = append(customProps, "ro.build.user=root")
	propLines = append(propLines, customProps...)
	newProp, err := os.Create(WorkPath + "build.prop")
	if err != nil {
		return err
	}
	os.Chmod(WorkPath+"build.prop", 0777)
	for _, line := range propLines {
		newProp.WriteString(line + "\n")
	}
	newProp.Sync()
	newProp.Close()
	vars.PatchLogger("Done with build.prop, moving on to other little files")
	vars.PatchLogger("/etc")
	os.WriteFile(WorkPath+"etc/timestamp", []byte(currentTime.Format("20060102150405")+"\n"), 0755)
	os.WriteFile(WorkPath+"etc/version", []byte(formattedTime+"\n"), 0755)
	issueContent := "msm-user " + formattedTime + " \\n \\l\n\n"
	os.WriteFile(WorkPath+"etc/issue", []byte(issueContent), 0755)
	issueNetContent := "msm-user " + formattedTime + " %h\n\n"
	os.WriteFile(WorkPath+"etc/issue.net", []byte(issueNetContent), 0755)
	osReleaseContent := `ID="msm"
NAME="msm"
VERSION="` + formattedTime + `"
VERSION_ID="` + formattedTime + `"
PRETTY_NAME="msm ` + formattedTime + `"
`
	os.WriteFile(WorkPath+"etc/os-release", []byte(osReleaseContent), 0755)
	os.WriteFile(WorkPath+"etc/os-version", []byte(version.Full+verAppendage+"\n"), 0755)
	os.WriteFile(WorkPath+"etc/os-version-base", []byte(version.Base+"\n"), 0755)
	os.WriteFile(WorkPath+"etc/os-version-code", []byte(version.Increment+"\n"), 0755)
	os.WriteFile(WorkPath+"etc/os-version-rev", []byte("wireos\n"), 0755)
	vars.PatchLogger("/anki/etc")
	cp.Copy(WorkPath+"anki/etc/version", WorkPath+"anki/etc/origversion")
	os.WriteFile(WorkPath+"anki/etc/version", []byte(version.Full+"\n"), 0755)
	os.WriteFile(WorkPath+"anki/etc/revision", []byte("wire-os\n"), 0755)
	return nil
}

func AddCorrectKernelModules(version vars.Version, target int) error {
	vars.PatchLogger("Replacing kernel modules with ones matching target: " + vars.Targets[target])
	moduleIf := "./resources/kernmods/" + vars.Targets[target] + "/modules"
	moduleOut := WorkPath + "usr/lib/modules"
	err := os.RemoveAll(moduleOut)
	if err != nil {
		return err
	}
	err = cp.Copy(moduleIf, moduleOut)
	if err != nil {
		return err
	}
	return nil
}

func MakeSysrootRW(version vars.Version, target int) error {
	var patchedFstabPathSuffix string
	if target == 3 {
		patchedFstabPathSuffix = "-orange"
	} else {
		patchedFstabPathSuffix = "-norm"
	}
	vars.PatchLogger("Copying fstab")
	os.Remove(WorkPath + "etc/fstab")
	err := cp.Copy("./resources/patches/MakeSysrootRW/fstab"+patchedFstabPathSuffix, WorkPath+"etc/fstab")
	if err != nil {
		return err
	}
	return nil
}

func PatchMountData(version vars.Version, target int) error {
	var patchedMountSuffix string
	if target == 3 {
		patchedMountSuffix = "-orange"
	} else {
		patchedMountSuffix = "-norm"
	}
	vars.PatchLogger("Copying mount-data")
	os.Remove(WorkPath + "etc/initscripts/mount-data")
	err := cp.Copy("./resources/patches/PatchMountData/mount-data"+patchedMountSuffix, WorkPath+"etc/initscripts/mount-data")
	if err != nil {
		return err
	}
	return nil
}

func HigherPriorityAnim(version vars.Version, target int) error {
	vars.PatchLogger("Copying new vic-anim service file")
	os.Remove(WorkPath + "lib/systemd/system/vic-anim.service")
	err := cp.Copy("./resources/patches/HigherPriorityAnim/vic-anim.service", WorkPath+"lib/systemd/system/vic-anim.service")
	if err != nil {
		return err
	}
	return nil
}

func AddSSHKey(version vars.Version, target int) error {
	if _, err := os.Stat(WorkPath + "etc/init.d/localsshuser.sh"); err == nil {
		os.Remove(WorkPath + "etc/init.d/localsshuser.sh")
		err := cp.Copy("./resources/patches/AddSSHKey/localsshuser.sh", WorkPath+"etc/init.d/localsshuser.sh")
		if err != nil {
			return err
		}
		exec.Command("/bin/bash", "-c", "chmod 0777 "+WorkPath+"etc/init.d/localsshuser.sh")
	}
	return nil
}
