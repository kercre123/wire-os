package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"

	manager "github.com/kercre123/wire-os/pkg/download-manager"
	handler "github.com/kercre123/wire-os/pkg/ota-handler"
	patcher "github.com/kercre123/wire-os/pkg/ota-patcher"
	"github.com/kercre123/wire-os/pkg/vars"
)

func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
}

func startWeb() {
	fmt.Println("Starting OTA fileserver at port 8080")
	http.Handle("/", http.FileServer(http.Dir("./out")))
	go http.ListenAndServe(":8080", nil)
}

func main() {
	if !isRoot() {
		fmt.Println("This program must be run as root, as it mounts filesystems.")
		os.Exit(0)
	}
	vars.InitDownloadManager()
	startWeb()
	for {
		fmt.Printf("\nEnter a command: ")
		var command string
		fmt.Scanln(&command)
		switch {
		case command == "download":
			fmt.Printf("\nEnter a URL to download: ")
			var URL string
			fmt.Scanln(&URL)
			manager.Download(URL)
		case command == "mount":
			fmt.Printf("\nEnter an OTA to mount: ")
			var OTA string
			fmt.Scanln(&OTA)
			version, err := manager.SplitVersion(OTA)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = handler.MountOTA(version)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case command == "unmount":
			fmt.Println("This will delete the image as well.")
			err := handler.UnmountImage("./work")
			if err != nil {
				fmt.Println(err)
				continue
			}
			os.Remove("./tmp/apq8009-robot-sysfs.img")
		case command == "pack":
			err := handler.PackOTA()
			if err != nil {
				fmt.Println(err)
				continue
			}
		case command == "delete":
			fmt.Printf("\nEnter an OTA version to delete: ")
			var OTA string
			fmt.Scanln(&OTA)
			version, err := manager.SplitVersion(OTA)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = manager.Delete(version)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Success!")
		case command == "patch":
			fmt.Printf("\nEnter desired output version (ex. 1.6.0.3331.ota): ")
			var OTA string
			fmt.Scanln(&OTA)
			version, err := manager.SplitVersion(OTA)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(vars.Targets)
			fmt.Printf("Enter desired target (0-3): ")
			var target string
			fmt.Scanln(&target)
			targetInd, err := strconv.Atoi(target)
			if err != nil || targetInd < 0 || targetInd > 3 {
				fmt.Println("target is invalid, it must be between 0 and 3")
				continue
			}
			fmt.Println("Running patches...")
			err = patcher.RunPatches(version, targetInd)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case command == "list":
			var downloaded []string
			var output []string
			storeEntries, _ := os.ReadDir("./store")
			for _, entry := range storeEntries {
				if strings.Contains(entry.Name(), ".ota") {
					downloaded = append(downloaded, strings.TrimSuffix(entry.Name(), ".ota"))
				}
			}
			outputEntries, _ := os.ReadDir("./out")
			for _, entry := range outputEntries {
				if strings.Contains(entry.Name(), ".ota") {
					output = append(output, entry.Name())
				}
			}
			fmt.Println("Downloaded OTAs: " + fmt.Sprint(downloaded))
			fmt.Println("Modified OTAs (accessible via webserver): " + fmt.Sprint(output))
		case command == "help":
			fmt.Println("(for now) These commands should be run without arguments. They are interactive")
			fmt.Println("download - Download an OTA from URL")
			fmt.Println("mount - Mount an OTA that has been downloaded")
			fmt.Println("unmount - Unmount currently-mounted OTA")
			fmt.Println("patch - Apply wireos patches to mounted OTA")
			fmt.Println("pack - Pack mounted OTA")
			fmt.Println("OTA format example: 1.6.0.3331")
			fmt.Println("Targets: 0 = dev, 1 = whiskey, 2 = oskr, 3 = orange")
		}
	}
}
