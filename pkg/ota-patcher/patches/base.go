package patches

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kercre123/wire-os/pkg/vars"
	cp "github.com/otiai10/copy"
)

const (
	WorkPath = "./work/"
)

var ProdServerConfig string = `{
	"jdocs": "jdocs.api.ddl.io:443",
	"tms": "token.api.ddl.io:443",
	"chipper": "chipper.api.ddl.io:443",
	"check": "conncheck.global.anki-services.com/ok",
	"logfiles": "s3://anki-device-logs-prod/victor",
	"appkey": "oDoa0quieSeir6goowai7f"
}
`

func AddVersion(version vars.Version, target int) error {
	// revision will eventually be hooked up to github commit
	vars.PatchLoggerName = "AddVersion"
	vars.PatchLogger("Modifying build.prop file")
	// read prop in OTA
	var propLines []string
	var customProps []string
	verAppendage := vars.Targets[target]
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
	customProps = append(customProps, "ro.build.version.release="+formattedTime)
	customProps = append(customProps, "ro.product.name=Vector")
	customProps = append(customProps, "ro.revision=wire_os")
	customProps = append(customProps, "ro.anki.version="+version.Full+verAppendage)
	customProps = append(customProps, "ro.anki.victor.version="+version.Full)
	customProps = append(customProps, "ro.build.fingerprint=v"+version.Full+"-wire_os"+version.Full+"-"+verAppendage+"-"+formattedTime)
	customProps = append(customProps, "ro.build.id=v"+version.Full+"-wire_os"+version.Full+"-"+verAppendage+"-"+formattedTime)
	customProps = append(customProps, "ro.build.display.id=v"+version.Full+verAppendage)
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
	vars.PatchLogger("anki/etc")
	os.WriteFile(WorkPath+"anki/etc/version", []byte(version.Full+"\n"), 0755)
	os.WriteFile(WorkPath+"anki/etc/revision", []byte("wire-os\n"), 0755)
	return nil
}

func AddCorrectKernelModules(version vars.Version, target int) error {
	vars.PatchLoggerName = "AddCorrectKernelModules"
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

func ProdServerEnv(version vars.Version, target int) error {
	vars.PatchLoggerName = "ProdServerEnv"
	vars.PatchLogger("Setting server env to prod")
	configPath := WorkPath + "anki/data/assets/cozmo_resources/config/server_config.json"
	_, err := os.ReadFile(configPath)
	if err != nil {
		vars.PatchLogger("Server config file does not exist. Not erroring because this is normal in older versions")
		return nil
	}
	err = os.WriteFile(configPath, []byte(ProdServerConfig), 0777)
	if err != nil {
		return err
	}
	return nil
}
