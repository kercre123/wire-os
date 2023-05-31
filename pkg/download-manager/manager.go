package manager

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"unicode"

	"github.com/kercre123/wire-os/pkg/vars"
	"github.com/schollz/progressbar/v3"
)

func SplitVersion(name string) (vars.Version, error) {
	// input: 1.9.0.1111.ota, or 0.11.19.ota
	version := vars.Version{}
	splitVersion := strings.Split(strings.TrimSuffix(name, ".ota"), ".")
	if len(splitVersion) == 3 {
		version.Base = splitVersion[0] + "." + splitVersion[1]
		version.Increment = RemoveLetters(splitVersion[2])
	} else if len(splitVersion) == 4 {
		version.Base = splitVersion[0] + "." + splitVersion[1] + "." + splitVersion[2]
		version.Increment = RemoveLetters(splitVersion[3])
	} else {
		return version, errors.New("incompatible version")
	}
	version.Full = strings.TrimSuffix(name, ".ota")
	return version, nil
}

func RemoveLetters(s string) string {
	// for OTAs which have "ep" or "oskr" at the end. the manifest stuff will deal with that
	var result string
	for _, r := range s {
		if unicode.IsDigit(r) {
			result += string(r)
		}
	}
	return result
}

func CreateDownloadEntry(name string, url string, completed bool) (id int) {
	entry := vars.DownloadInfo{}
	// error already checked
	version, _ := SplitVersion(name)
	entry.Version = version
	entry.URL = url
	entry.FileName = name
	entry.Completed = completed
	vars.DownloadsInfo = append(vars.DownloadsInfo, entry)
	vars.WriteJson(vars.DownloadsInfo, vars.DownloadsJson)
	return len(vars.DownloadsInfo) - 1
}

func Download(URL string) {
	splitURL := strings.Split(URL, "/")
	filename := splitURL[len(splitURL)-1]
	_, err := SplitVersion(filename)
	if err != nil {
		fmt.Println("OTA filename does not match standard (x.x.x.xxxx or x.x.xxxx)")
		return
	}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		fmt.Println("Error downloading OTA")
		fmt.Println(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error downloading OTA")
		fmt.Println(err)
	}
	defer resp.Body.Close()

	f, _ := os.Create("./store/" + filename)
	defer f.Close()
	downloadID := CreateDownloadEntry(filename, URL, false)
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		filename,
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)
	vars.DownloadsInfo[downloadID].Completed = true
	vars.WriteJson(vars.DownloadsInfo, vars.DownloadsJson)
	fmt.Println("Download of " + filename + " has completed!")
}
