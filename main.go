package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	manager "github.com/kercre123/wire-os/pkg/download-manager"
	handler "github.com/kercre123/wire-os/pkg/ota-handler"
	"github.com/kercre123/wire-os/pkg/vars"
)

func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
}

func main() {
	if !isRoot() {
		fmt.Println("This program must be run as root, as it mounts filesystems.")
		os.Exit(0)
	}
	vars.InitDownloadManager()
	for {
		fmt.Printf("\nEnter a command: ")
		var command string
		fmt.Scanln(&command)
		switch {
		case command == "download":
			manager.Download("http://ota.global.anki-services.com/vic/prod/full/1.6.0.3331.ota")
		case command == "mount":
			handler.MountOTA("./store/1.6.0.3331.ota")
		case command == "pack":
			err := handler.PackOTA(vars.Version{
				Base:      "1.7.0",
				Increment: "1",
				Full:      "1.7.0.1",
			}, 0)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
