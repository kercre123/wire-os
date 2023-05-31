package vars

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const (
	JsonPath      = "./json/"
	DownloadsJson = JsonPath + "downloads.json"
)

var Targets []string = []string{"dev", "whiskey", "oskr", "orange"}

type Modifier struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
	ID          uint32 `json:"id"`
}

type Version struct {
	Base      string `json:"base"`
	Increment string `json:"increment"`
	Full      string `json:"full"`
}

type DownloadInfo struct {
	Version   Version    `json:"version"`
	FileName  string     `json:"filename"`
	URL       string     `json:"url"`
	Completed bool       `json:"completed"`
	Modifiers []Modifier `json:"modifiers"`
}

// used for runtime, gets loaded in at start of program
var DownloadsInfo []DownloadInfo

func InitDownloadManager() {
	downloads, err := os.Open(DownloadsJson)
	if err != nil {
		fmt.Println("Downloads JSON does not exist")
		return
	}
	jsonBytes, err := io.ReadAll(downloads)
	if err != nil {
		fmt.Println("Error reading downloads JSON")
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(jsonBytes, &DownloadsInfo)
	if err != nil {
		fmt.Println("Error unmarshaling downloads JSON")
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully loaded downloads JSON, listing:")
	fmt.Println(DownloadsInfo)
}

func WriteJson(data any, filepath string) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Println("Writing to " + filepath)
	err = os.WriteFile(filepath, jsonBytes, 0777)
	if err != nil {
		return err
	}
	return nil
}
