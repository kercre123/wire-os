package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	manager "github.com/kercre123/wire-os/pkg/download-manager"
	handler "github.com/kercre123/wire-os/pkg/ota-handler"
	patcher "github.com/kercre123/wire-os/pkg/ota-patcher"
	"github.com/kercre123/wire-os/pkg/vars"
)

// func isRoot() bool {
// 	currentUser, err := user.Current()
// 	if err != nil {
// 		log.Fatalf("[isRoot] Unable to get current user: %s", err)
// 	}
// 	return currentUser.Username == "root"
// }

func startWeb() {
	fmt.Println("Starting OTA fileserver at port 8080")
	http.Handle("/", http.FileServer(http.Dir("./out")))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}

func main() {
	// if !isRoot() {
	// 	fmt.Println("This program must be run as root, as it mounts filesystems.")
	// 	os.Exit(1)
	// }

	vars.InitDownloadManager()

	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	downloadURL := downloadCmd.String("url", "", "URL to download the OTA from")

	mountCmd := flag.NewFlagSet("mount", flag.ExitOnError)
	mountOTA := mountCmd.String("ota", "", "OTA version to mount")

	unmountCmd := flag.NewFlagSet("unmount", flag.ExitOnError)

	packCmd := flag.NewFlagSet("pack", flag.ExitOnError)

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteOTA := deleteCmd.String("ota", "", "OTA version to delete")

	patchCmd := flag.NewFlagSet("patch", flag.ExitOnError)
	patchOutput := patchCmd.String("output", "", "Desired output version (e.g., 1.6.0.3331.ota)")
	patchTarget := patchCmd.Int("target", -1, "Desired target (0-3)")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	startWebCmd := flag.NewFlagSet("startweb", flag.ExitOnError)

	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Expected 'download', 'mount', 'unmount', 'pack', 'delete', 'patch', 'list', 'startweb' or 'help' subcommands.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "download":
		downloadCmd.Parse(os.Args[2:])
		if *downloadURL == "" {
			fmt.Println("Usage: download -url=<URL>")
			os.Exit(1)
		}
		manager.Download(*downloadURL)

	case "mount":
		mountCmd.Parse(os.Args[2:])
		if *mountOTA == "" {
			fmt.Println("Usage: mount -ota=<OTA version>")
			os.Exit(1)
		}
		version, err := manager.SplitVersion(*mountOTA)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		err = handler.MountOTA(version)
		if err != nil {
			fmt.Printf("Error mounting OTA: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("OTA mounted successfully.")

	case "unmount":
		unmountCmd.Parse(os.Args[2:])
		fmt.Println("This will delete the image as well.")
		err := handler.UnmountImage("./work")
		if err != nil {
			fmt.Printf("Error unmounting OTA: %v\n", err)
			os.Exit(1)
		}
		err = os.Remove("./tmp/apq8009-robot-sysfs.img")
		if err != nil {
			fmt.Printf("Error removing image: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("OTA unmounted and image deleted successfully.")

	case "pack":
		packCmd.Parse(os.Args[2:])
		err := handler.PackOTA()
		if err != nil {
			fmt.Printf("Error packing OTA: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("OTA packed successfully.")

	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if *deleteOTA == "" {
			fmt.Println("Usage: delete -ota=<OTA version>")
			os.Exit(1)
		}
		version, err := manager.SplitVersion(*deleteOTA)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		err = manager.Delete(version)
		if err != nil {
			fmt.Printf("Error deleting OTA: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("OTA deleted successfully.")

	case "patch":
		patchCmd.Parse(os.Args[2:])
		if *patchOutput == "" || *patchTarget == -1 {
			fmt.Println("Usage: patch -output=<output OTA version> -target=<0-3>")
			os.Exit(1)
		}
		version, err := manager.SplitVersion(*patchOutput)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if *patchTarget < 0 || *patchTarget > 3 {
			fmt.Println("Error: target must be between 0 and 3")
			os.Exit(1)
		}
		fmt.Println("Running patches...")
		err = patcher.RunPatches(version, *patchTarget)
		if err != nil {
			fmt.Printf("Error applying patches: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Patches applied successfully.")

	case "list":
		listCmd.Parse(os.Args[2:])
		var downloaded []string
		var output []string
		storeEntries, err := os.ReadDir("./store")
		if err != nil {
			fmt.Printf("Error reading store directory: %v\n", err)
			os.Exit(1)
		}
		for _, entry := range storeEntries {
			if strings.Contains(entry.Name(), ".ota") {
				downloaded = append(downloaded, strings.TrimSuffix(entry.Name(), ".ota"))
			}
		}
		outputEntries, err := os.ReadDir("./out")
		if err != nil {
			fmt.Printf("Error reading out directory: %v\n", err)
			os.Exit(1)
		}
		for _, entry := range outputEntries {
			if strings.Contains(entry.Name(), ".ota") {
				output = append(output, entry.Name())
			}
		}
		fmt.Println("Downloaded OTAs: " + fmt.Sprint(downloaded))
		fmt.Println("Modified OTAs (accessible via webserver): " + fmt.Sprint(output))

	case "startweb":
		startWebCmd.Parse(os.Args[2:])
		startWeb()

	case "help":
		helpCmd.Parse(os.Args[2:])
		fmt.Println("Available subcommands:")
		fmt.Println("  download -url=<URL>          Download an OTA from the specified URL.")
		fmt.Println("  mount -ota=<OTA version>     Mount a downloaded OTA.")
		fmt.Println("  unmount                      Unmount the currently mounted OTA and delete the image.")
		fmt.Println("  pack                         Pack the currently mounted OTA.")
		fmt.Println("  delete -ota=<OTA version>    Delete a specified OTA version.")
		fmt.Println("  patch -output=<OTA> -target=<0-3>  Apply patches to the specified OTA with the given target.")
		fmt.Println("  list                         List downloaded and modified OTAs.")
		fmt.Println("  startweb                     Start the OTA fileserver on port 8080.")
		fmt.Println("  help                         Show this help message.")
		fmt.Println("\nOTA format example: 1.6.0.3331")
		fmt.Println("Targets:")
		fmt.Println("  0 = dev")
		fmt.Println("  1 = whiskey")
		fmt.Println("  2 = oskr")
		fmt.Println("  3 = orange")

	default:
		fmt.Printf("Unknown subcommand: %s\n", os.Args[1])
		fmt.Println("Available subcommands: download, mount, unmount, pack, delete, patch, list, startweb, help")
		os.Exit(1)
	}
}
