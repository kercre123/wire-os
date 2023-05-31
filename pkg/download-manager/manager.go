package manager

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"unicode"

	"github.com/kercre123/wire-os/pkg/vars"
	"github.com/schollz/progressbar/v3"
	"gopkg.in/ini.v1"
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
	} else if strings.Contains(name, "latest") {
		version.Base = "0.0.0"
		version.Increment = "0000"
		version.Full = "0.0.0.0000"
		return version, errors.New("latest.ota")
	} else {
		return version, errors.New("incompatible version")
	}
	version.Full = version.Base + "." + version.Increment
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
	getVerFromManifest := false
	_, err := SplitVersion(filename)
	if err != nil {
		if err.Error() == "latest.ota" {
			getVerFromManifest = true
		} else {
			fmt.Println("OTA filename does not match standard (x.x.x.xxxx or x.x.xxxx)")
			return
		}
	}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		fmt.Println("Error downloading OTA")
		fmt.Println(err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error downloading OTA")
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	var manifestBuffer bytes.Buffer
	if getVerFromManifest {
		reader := tar.NewReader(io.TeeReader(resp.Body, &manifestBuffer))
		for {
			head, err := reader.Next()
			if err != nil {
				fmt.Println("Error downloading OTA (when getting from manifest)")
				fmt.Println(err)
				return
			}
			if strings.Contains(head.Name, "manifest.ini") {
				manifestBytes, _ := io.ReadAll(reader)
				manifest, err := ini.Load(manifestBytes)
				if err != nil {
					fmt.Println("Error downloading OTA (when getting from manifest)")
					fmt.Println(err)
					return
				}
				meta, err := manifest.GetSection("META")
				if err != nil {
					fmt.Println("Error downloading OTA (when getting from manifest)")
					fmt.Println(err)
					return
				}
				key, err := meta.GetKey("update_version")
				if err != nil {
					fmt.Println("Error downloading OTA (when getting from manifest)")
					fmt.Println(err)
					return
				}
				version, _ := SplitVersion(key.String() + ".ota")
				filename = version.Full + ".ota"
				break
			}
		}
	}

	f, _ := os.Create("./store/" + filename)
	defer f.Close()
	downloadID := CreateDownloadEntry(filename, URL, false)
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		filename,
	)
	if getVerFromManifest {
		io.Copy(io.MultiWriter(f, bar), io.MultiReader(&manifestBuffer, resp.Body))
	} else {
		io.Copy(io.MultiWriter(f, bar), resp.Body)
	}
	vars.DownloadsInfo[downloadID].Completed = true
	vars.WriteJson(vars.DownloadsInfo, vars.DownloadsJson)
	fmt.Println("Download of " + filename + " has completed!")
}

func Delete(version vars.Version) error {
	matched := false
	var newInfo []vars.DownloadInfo
	for _, info := range vars.DownloadsInfo {
		if info.Version == version {
			fmt.Println("Deleting " + version.Full + ".ota")
			os.Remove("./store/" + info.FileName)
			matched = true
			for _, info := range vars.DownloadsInfo {
				if info.Version != version {
					newInfo = append(newInfo, info)
				}
			}
		}
	}
	if matched {
		vars.DownloadsInfo = newInfo
		vars.WriteJson(vars.DownloadsInfo, vars.DownloadsJson)
	} else {
		return errors.New("file not in downloads json")
	}
	return nil
}
