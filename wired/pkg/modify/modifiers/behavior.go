package modifiers

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/kercre123/wire-os/wired/pkg/vars"
	cp "github.com/otiai10/copy"
)

type AG struct {
	Animations []struct {
		Name            string  `json:"Name"`
		Weight          float64 `json:"Weight"`
		CooldownTimeSec float64 `json:"CooldownTime_Sec"`
		Mood            string  `json:"Mood"`
	} `json:"Animations"`
}

func RestartAnkiServices() {
	fmt.Println("Restarting anki processes (screen will go blank for a sec)...")
	exec.Command("/bin/systemctl", "stop", "anki-robot.target").Run()
	time.Sleep(time.Second * 3)
	exec.Command("/bin/systemctl", "start", "anki-robot.target").Run()
	fmt.Println("Anki processes starting back up!")
}

func NoSnore_Apply() error {
	ModPath := vars.ModifiersEFiles + "NoSnore/"
	snoreAGPath := vars.AnkiResources + "assets/animationGroups/gotoSleep/ag_gotosleep_sleeping.json"
	var snoreAG AG
	os.MkdirAll(ModPath, 0777)
	// back up file for future removing
	cp.Copy(snoreAGPath, ModPath+"ag_gotosleep_sleeping.json")
	snoreAGBytes, err := os.ReadFile(snoreAGPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(snoreAGBytes, &snoreAG)
	if err != nil {
		return err
	}
	for ind, ag := range snoreAG.Animations {
		if ag.Weight > 0.0 && ag.Name != "anim_gotosleep_sleeping_02" {
			snoreAG.Animations[ind].Weight = 0.0
		}
	}
	jsonBytes, err := json.Marshal(snoreAG)
	if err != nil {
		return err
	}
	err = os.WriteFile(snoreAGPath, jsonBytes, 0777)
	if err != nil {
		return err
	}
	RestartAnkiServices()
	return nil
}

func NoSnore_Remove() error {
	ModPath := vars.ModifiersEFiles + "NoSnore/"
	snoreAGPath := vars.AnkiResources + "assets/animationGroups/gotoSleep/ag_gotosleep_sleeping.json"
	cp.Copy(ModPath+"ag_gotosleep_sleeping.json", snoreAGPath)
	RestartAnkiServices()
	return nil
}
