package mods

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	ba "github.com/kercre123/wire-os/wired/raw"
	"github.com/kercre123/wire-os/wired/vars"
	cp "github.com/otiai10/copy"
)

type BootAnim struct {
	vars.Modification
}

func NewBootAnim() *BootAnim {
	return &BootAnim{}
}

var BootAnim_Current BootAnim_AcceptJSON

type BootAnim_AcceptJSON struct {
	Default bool   `json:"default"`
	GifData string `json:"gifdata"`
}

func (modu *BootAnim) Name() string {
	return "BootAnim"
}

func (modu *BootAnim) Description() string {
	return "Boot animation from GIF."
}

func (modu *BootAnim) RestartRequired() bool {
	return false
}

func (modu *BootAnim) DefaultJSON() any {
	return BootAnim_AcceptJSON{
		Default: true,
	}
}

func BootAnim_HTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/mods/custom/TestBootAnim" {
		BootAnim_Show()
		vars.HTTPSuccess(w, r)
	} else {
		vars.HTTPError(w, r, "404 not found")
	}
}

func BootAnim_Show() {
	// show anim on screen for 10 seconds
	cmd := exec.Command("/bin/bash", "-c", "/anki/bin/vic-bootAnim")
	vars.StopVic()
	go func() {
		cmd.Run()
	}()
	time.Sleep(time.Second * 15)
	cmd.Process.Kill()
	time.Sleep(time.Second * 1)
	vars.StartVic()
}

func (modu *BootAnim) Save(where string, in string) error {
	moduin, ok := modu.DefaultJSON().(BootAnim_AcceptJSON)
	if !ok {
		return errors.New("internal mod error: Save(), DefaultJSON not correct type")
	}
	json.Unmarshal([]byte(in), &moduin)
	saveJson, err := json.Marshal(moduin)
	if err != nil {
		return err
	}
	os.MkdirAll(vars.GetModDir(modu, where), 0777)
	os.WriteFile(vars.GetModDir(modu, where)+"/saved.json", saveJson, 0777)
	return nil
}

func (modu *BootAnim) Load() error {
	moduin, ok := modu.DefaultJSON().(BootAnim_AcceptJSON)
	if !ok {
		return errors.New("internal mod error: Load(), DefaultJSON not correct type")
	}
	file, err := os.ReadFile(vars.GetModDir(modu, "/") + "saved.json")
	if err != nil {
		defaultJson, _ := json.Marshal(moduin)
		modu.Do("/", string(defaultJson))
		return nil
	}
	json.Unmarshal(file, &moduin)
	BootAnim_Current = moduin
	return nil
}

func (modu *BootAnim) Accepts() string {
	str, ok := modu.DefaultJSON().(BootAnim_AcceptJSON)
	if !ok {
		log.Fatal("BootAnim Accepts(): not correct type")
	}
	marshedJson, err := json.Marshal(str)
	if err != nil {
		log.Fatal(err)
	}
	return string(marshedJson)
}

func (modu *BootAnim) Current() string {
	marshalled, _ := json.Marshal(BootAnim_Current)
	return string(marshalled)
}

func (modu *BootAnim) Do(where string, in string) error {
	moduin, ok := modu.DefaultJSON().(BootAnim_AcceptJSON)
	if !ok {
		return errors.New("internal mod error: Do(), DefaultJSON not correct type")
	}
	err := json.Unmarshal([]byte(in), &moduin)
	if err != nil {
		return err
	}
	var reqIn BootAnim_AcceptJSON
	err = json.Unmarshal([]byte(in), &reqIn)
	if err != nil {
		return err
	}
	if reqIn.Default {
		os.Remove(where + vars.VectorResources + "config/engine/animations/boot_anim.raw")
		cp.Copy(vars.GetModDir(modu, where)+"orig_boot_anim.raw", where+vars.VectorResources+"config/engine/animations/boot_anim.raw")
	} else {
		gifBytes, err := base64.StdEncoding.DecodeString(reqIn.GifData)
		if err != nil {
			return err
		}
		os.Remove(where + "data/boot_anim_new.raw")
		err = ba.GifToBootAnimation(gifBytes, where+"data/boot_anim_new.raw")
		if err != nil {
			os.Remove(where + "data/boot_anim_new.raw")
			return err
		}
		os.Remove(where + vars.VectorResources + "config/engine/animations/boot_anim.raw")
		cp.Copy(where+"data/boot_anim_new.raw", where+vars.VectorResources+"config/engine/animations/boot_anim.raw")
		os.Remove(where + "data/boot_anim_new.raw")
	}
	modu.Save(where, in)
	BootAnim_Current = moduin
	return nil
}
