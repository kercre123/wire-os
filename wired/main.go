package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kercre123/wire-os/wired/mods"
	"github.com/kercre123/wire-os/wired/vars"
)

var EnabledMods []vars.Modification = []vars.Modification{
	mods.NewFreqChange(),
	mods.NewRainbowLights(),
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

func ModHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/api/mods/modify/"):
		modFromURL := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/mods/modify/"))
		jsonFromReq, err := io.ReadAll(r.Body)
		if err != nil {
			HTTPError(w, r, "error reading request body: "+err.Error())
			return
		}
		mod, err := vars.FindMod(modFromURL)
		if err != nil {
			HTTPError(w, r, "mod does not exist")
			return
		}
		err = mod.Do(vars.SysRoot, string(jsonFromReq))
		if err != nil {
			HTTPError(w, r, "error running mod: "+err.Error())
			return
		}
		HTTPSuccess(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/mods/current/"):
		modFromURL := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/mods/current/"))
		mod, err := vars.FindMod(modFromURL)
		if err != nil {
			HTTPError(w, r, "mod does not exist")
			return
		}
		fmt.Fprint(w, mod.Current())
	case strings.HasPrefix(r.URL.Path, "/api/mods/needsrestart/"):
		modFromURL := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/mods/needsrestart/"))
		mod, err := vars.FindMod(modFromURL)
		if err != nil {
			HTTPError(w, r, "mod does not exist")
			return
		}
		fmt.Fprint(w, mod.RestartRequired())
	case strings.HasPrefix(r.URL.Path, "/api/restartvic"):
		vars.RestartVic()
		HTTPSuccess(w, r)
	default:
		HTTPError(w, r, "404 not found")
	}
}

func main() {
	vars.EnabledMods = EnabledMods
	vars.InitMods()
	startweb()
}

func startweb() {
	fmt.Println("starting web at port 8080")
	http.Handle("/", http.FileServer(http.Dir("/etc/wired/webroot")))
	http.HandleFunc("/api/mods/modify/", ModHandler)
	http.HandleFunc("/api/mods/current/", ModHandler)
	http.HandleFunc("/api/mods/needsrestart/", ModHandler)
	http.HandleFunc("/api/restartvic", ModHandler)
	http.ListenAndServe(":8080", nil)
}

func startshell() {
	fmt.Println("starting shell")
	for {
		fmt.Printf("\n#~ ")
		var in string
		fmt.Scanln(&in)
		switch {
		case in == "list":
			for _, mod := range EnabledMods {
				fmt.Println("")
				fmt.Println(mod.Name())
				fmt.Println(mod.Description())
				fmt.Println(mod.Accepts())
				fmt.Println("")
			}
		case in == "find":
			var find string
			fmt.Println("which mod would you like to find? ")
			fmt.Printf("Name: ")
			fmt.Scanln(&find)
			mod, err := vars.FindMod(find)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(mod.Name() + " found")
			fmt.Println("Enter a JSON input for Do()")
			fmt.Printf("Input: ")
			var doinput string
			fmt.Scanln(&doinput)
			err = mod.Do("/", doinput)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}
