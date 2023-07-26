package vars

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	SysRoot = "/"
	// for applying mods to updates
	Update_SysRoot  = "/mnt/"
	VectorResources = "anki/data/assets/cozmo_resources/"
)

type Modification interface {
	Name() string
	Description() string
	Accepts() string
	DefaultJSON() any
	Save(string, string) error
	// note: Load() runs at init of program
	Load() error
	// current settings of mod
	Current() string
	// fs root
	ToFS(string)
	RestartRequired() bool
	Do(string, string) error
}

type BaseModification struct {
	Modification
	ModName            string
	ModDescription     string
	VicRestartRequired bool
}

func (bc *BaseModification) Name() string {
	return bc.ModName
}

func (bc *BaseModification) Description() string {
	return bc.ModDescription
}

func (bc *BaseModification) RestartRequired() bool {
	return bc.VicRestartRequired
}

var EnabledMods []Modification

func GetModDir(mod Modification, where string) string {
	return where + "etc/wired/mods/" + mod.Name() + "/"
	//return "./modtest/" + mod.Name() + "/"
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
		fmt.Println("Loading " + mod.Name() + "...")
		mod.Load()
	}
}

func StopVic() {
	Behavior("DevBaseBehavior")
	time.Sleep(time.Second * 1)
	exec.Command("/bin/bash", "-c", "systemctl stop anki-robot.target").Output()
	time.Sleep(time.Second * 4)
}

func StartVic() {
	exec.Command("/bin/bash", "-c", "systemctl start anki-robot.target").Output()
	time.Sleep(time.Second * 3)
}

func RestartVic() {
	StopVic()
	StartVic()
}

type HTTPStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func HTTPSuccess(w http.ResponseWriter, r *http.Request) {
	var status HTTPStatus
	status.Status = "success"
	successBytes, _ := json.Marshal(status)
	w.Write(successBytes)
}

func HTTPError(w http.ResponseWriter, r *http.Request, err string) {
	var status HTTPStatus
	status.Status = "error"
	status.Message = err
	errorBytes, _ := json.Marshal(status)
	w.Write(errorBytes)
}

type BehaviorMessage struct {
	Type   string `json:"type"`
	Module string `json:"module"`
	Data   struct {
		BehaviorName     string `json:"behaviorName"`
		PresetConditions bool   `json:"presetConditions"`
	} `json:"data"`
}

//{"type":"data","module":"behaviors","data":{"behaviorName":"DevBaseBehavior","presetConditions":false}}

func Behavior(behavior string) {
	u := url.URL{Scheme: "ws", Host: "localhost:8888", Path: "/socket"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	message := BehaviorMessage{
		Type:   "data",
		Module: "behaviors",
		Data: struct {
			BehaviorName     string `json:"behaviorName"`
			PresetConditions bool   `json:"presetConditions"`
		}{
			BehaviorName:     behavior,
			PresetConditions: false,
		},
	}

	marshaledMessage, err := json.Marshal(message)
	if err != nil {
		log.Fatal("marshal:", err)
	}

	err = c.WriteMessage(websocket.TextMessage, marshaledMessage)
	if err != nil {
		log.Fatal("write:", err)
	}
}
