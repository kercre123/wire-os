package vars

import (
	"encoding/json"
	"fmt"
	"os"
)

// things which need to be done:
// 1. OnUpdate func(bool) error
// 	- wireos-modifiers.json should get saved to /data as well, so that patches can get added upon OS update
// 2. the actual web ui
// 3. maybe a system to wake the robot up before patching, so that restarting the anki processes is more reliable

const (
	ModifiersJSON   = "/etc/wireos-modifiers.json"
	ModifiersEFiles = "/etc/wireos-modifiers"
	AnkiResources   = "/anki/data/assets/cozmo_resources/"
)

var Modifiers []Modifier
var LoadedModifiers []LoadedModifier

type Modifier struct {
	Name        string
	Description string
	Apply       func() error
	Remove      func() error
	HasInitFunc bool
	Init        func(bool) error
}

type LoadedModifier struct {
	// matches up to index in array of Modifiers
	ModifierID int  `json:"id"`
	Applied    bool `json:"applied"`
}

type ReportedModifier struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
	ModifierID  int    `json:"id"`
	Applied     bool   `json:"applied"`
}

func LoadModifiersJSON() {
	jsonBytes, err := os.ReadFile(ModifiersJSON)
	if err != nil {
		fmt.Println(ModifiersJSON + " does not exist, will create")
		return
	}
	err = json.Unmarshal(jsonBytes, &LoadedModifiers)
	if err != nil {
		fmt.Println("Error unmarshalling modifiers JSON, will create new file")
		return
	}
}

func SaveModifiersJSON() {
	fmt.Println("Writing modifiers JSON (" + ModifiersJSON + ")")
	newJson, err := json.Marshal(LoadedModifiers)
	if err != nil {
		fmt.Println("Error marshalling JSON: " + err.Error())
		return
	}
	err = os.WriteFile(ModifiersJSON, newJson, 0777)
	if err != nil {
		fmt.Println("Error writing modifiers JSON: " + err.Error())
		return
	}
}

func RunInitFuncs() {
	for _, lmod := range LoadedModifiers {
		if Modifiers[lmod.ModifierID].HasInitFunc {
			fmt.Println("Running init function for " + Modifiers[lmod.ModifierID].Name)
			Modifiers[lmod.ModifierID].Init(lmod.Applied)
		}
	}
}

func Init() {
	LoadModifiersJSON()
	RunInitFuncs()
}
