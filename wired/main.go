package main

import (
	"fmt"
	"net/http"

	"github.com/kercre123/wire-os/wired/mods"
	"github.com/kercre123/wire-os/wired/vars"
)

var EnabledMods []vars.Modification = []vars.Modification{
	mods.NewFreqChange(),
	mods.NewWakeWord(),
}

func main() {
	vars.EnabledMods = EnabledMods
	vars.InitMods()
	startweb()
}

func startweb() {
	fmt.Println("starting web at port 8080")
	http.Handle("/", http.FileServer(http.Dir("/etc/wired/webroot")))
	mods.ImplHTTP()
	http.ListenAndServe(":8080", nil)
}
