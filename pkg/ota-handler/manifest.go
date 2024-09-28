package handler

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/kercre123/wire-os/pkg/vars"
	"gopkg.in/ini.v1"
)

/*
example ini:
[META]
manifest_version=0.9.2
update_version=2.0.0.6074
ankidev=0
num_images=2
reboot_after_install=0
[BOOT]
encryption=1
delta=0
compression=gz
wbits=31
bytes=12869632
sha256=c45087ff6b0f581b5d42e9616e14beced16ea3928f01e950aa9fea5aa4b524c5
[SYSTEM]
encryption=1
delta=0
compression=gz
wbits=31
bytes=608743424
sha256=43ee23615df81a7888822b596c98d0fa52b231f47ec0a4abf7e2de0c9d92c8c3
*/

type BootInfo struct {
	Filename string `json:"filename"`
	Length   int    `json:"length"`
	Hash     string `json:"hash"`
}

var BootInfos []BootInfo

func FindBootLength(target int) string {
	targetString := vars.Targets[target]
	bootLengthJsonFile, _ := os.ReadFile("./json/bootLengths.json")
	json.Unmarshal(bootLengthJsonFile, &BootInfos)
	for _, info := range BootInfos {
		if info.Filename == "./boots/"+targetString+".img.gz" {
			return strconv.Itoa(info.Length)
		}
	}
	return "0"
}

func FindBootHash(target int) string {
	targetString := vars.Targets[target]
	bootLengthJsonFile, _ := os.ReadFile("./json/bootLengths.json")
	json.Unmarshal(bootLengthJsonFile, &BootInfos)
	for _, info := range BootInfos {
		if info.Filename == "./boots/"+targetString+".img.gz" {
			return info.Hash
		}
	}
	return "0"
}

func FindImageLength(image string) string {
	file, _ := os.ReadFile(image)
	return strconv.Itoa(len(file))
}

func FindImageHash(image string) string {
	filename := filepath.Join(os.TempDir(), "apq8009-robot-sysfs.img")
	file, _ := os.ReadFile(image)
	fmt.Println("[FindImageHash] Compressing...")
	compressed, _ := CompressBytes(file)
	fmt.Println("[FindImageHash] Decompressing...")
	decompressed, _ := DecompressBytes(compressed)
	os.WriteFile(filename, decompressed, 0777)
	fmt.Println("[FindImageHash] Finding hash of system image...")
	finalfile, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file: " + err.Error())
		return ""
	}
	h := sha256.New()
	if _, err := io.Copy(h, finalfile); err != nil {
		fmt.Println("error getting hash: " + err.Error())
		return ""
	}
	os.Remove(filename)
	return hex.EncodeToString(h.Sum(nil))
}

func CreateManifest(version vars.Version, image string, target int) []byte {
	//imageBytes, _ := os.ReadFile(image)
	var buf []byte
	buffer := bytes.NewBuffer(buf)
	manifest := ini.Empty()
	meta, _ := manifest.NewSection("META")
	meta.NewKey("manifest_version", "0.9.2")
	meta.NewKey("update_version", version.Full)
	meta.NewKey("ankidev", "1")
	meta.NewKey("num_images", "2")
	meta.NewKey("reboot_after_install", "0")
	boot, _ := manifest.NewSection("BOOT")
	boot.NewKey("encryption", "1")
	boot.NewKey("delta", "0")
	boot.NewKey("compression", "gz")
	boot.NewKey("wbits", "31")
	boot.NewKey("bytes", FindBootLength(target))
	boot.NewKey("sha256", FindBootHash(target))
	sysfs, _ := manifest.NewSection("SYSTEM")
	sysfs.NewKey("encryption", "1")
	sysfs.NewKey("delta", "0")
	sysfs.NewKey("compression", "gz")
	sysfs.NewKey("wbits", "31")
	sysfs.NewKey("bytes", FindImageLength(image))
	sysfs.NewKey("sha256", FindImageHash(image))
	manifest.WriteTo(buffer)
	return buffer.Bytes()
}
