package handler

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	manager "github.com/kercre123/wire-os/pkg/download-manager"
	"github.com/kercre123/wire-os/pkg/vars"
)

func EncryptBytes(in []byte) ([]byte, error) {
	inFileName := os.TempDir() + "/wireos.gz"
	outFileName := os.TempDir() + "/wireos-encrypted.gz"
	err := os.WriteFile(inFileName, in, 0777)
	if err != nil {
		return in, err
	}
	err = exec.Command("openssl", "enc", "-e", "-aes-256-ctr", "-pass", "file:"+PassphraseFile, "-md", "md5", "-in", inFileName, "-out", outFileName).Run()
	if err != nil {
		return in, err
	}
	decrypted, err := os.ReadFile(outFileName)
	if err != nil {
		return in, err
	}
	os.Remove(inFileName)
	os.Remove(outFileName)
	return decrypted, nil
}

func CompressBytes(in []byte) ([]byte, error) {
	var buf []byte
	writer := bytes.NewBuffer(buf)
	gzwriter := gzip.NewWriter(writer)
	gzwriter.Name = "apq8009-robot-sysfs.img.gz"
	_, err := gzwriter.Write(in)
	if err != nil {
		return []byte{}, err
	}
	gzwriter.Close()
	return writer.Bytes(), nil
}

func UnmountImage(path string) error {
	// err := exec.Command("guestunmount", path).Run()
	// if err != nil {
	// 	fmt.Println("couldn't guestunmount?")
	// 	return err
	// }
	err := exec.Command("umount", path).Run()
	if err != nil {
		fmt.Println("couldn't unmount?")
		return err
	}
	exec.Command("sync").Run()
	err = exec.Command("fsck.ext4", "-f", "-y", "./tmp/apq8009-robot-sysfs.img").Run()
	if err != nil {
		fmt.Println("first fsck failed, running again")
		out, err := exec.Command("fsck.ext4", "-f", "-y", "./tmp/apq8009-robot-sysfs.img").Output()
		if err != nil {
			fmt.Println(string(out))
			return err
		}
		fmt.Println("second fsck succeeded!")
	}
	return nil
}

func PackTar(systembytes []byte, target int, manifest []byte, version vars.Version) error {
	// 0 = dev, 1 = whiskey, 2 = oskr, 3 = orange(fills a very specific niche of bot)
	var buf []byte
	if target > 4 || target < 0 {
		return errors.New("target out of range, must be between 1 and 4")
	}
	bootbytes, err := os.ReadFile("./resources/boots/" + vars.Targets[target] + ".img.gz")
	if err != nil {
		return err
	}
	file := bytes.NewBuffer(buf)
	tarFile := tar.NewWriter(file)
	// write manifest
	err = tarFile.WriteHeader(&tar.Header{
		Name: "manifest.ini",
		Size: int64(len(manifest)),
		Mode: 0777,
	})
	if err != nil {
		return err
	}
	tarFile.Write(manifest)
	// sha would go here
	err = tarFile.WriteHeader(&tar.Header{
		Name: "apq8009-robot-boot.img.gz",
		Size: int64(len(bootbytes)),
		Mode: 0777,
	})
	if err != nil {
		return err
	}
	tarFile.Write(bootbytes)
	err = tarFile.WriteHeader(&tar.Header{
		Name: "apq8009-robot-sysfs.img.gz",
		Size: int64(len(systembytes)),
		Mode: 0777,
	})
	if err != nil {
		return err
	}
	tarFile.Write(systembytes)
	err = tarFile.Close()
	if err != nil {
		return err
	}
	err = os.WriteFile("./out/"+version.Full+".ota", file.Bytes(), 0777)
	if err != nil {
		return err
	}
	os.Remove("./tmp/apq8009-robot-sysfs.img")
	return nil
}

// version and target will be picked out from OTA contents
func PackOTA() error {
	var target int
	if !vars.IsMounted() {
		err := MountImage("./tmp/apq8009-robot-sysfs.img", "./work")
		if err != nil {
			return errors.New("image wasn't mounted, failure upon trying to mount")
		}
	}
	// get version from build.prop
	prop, err := os.Open("./work/build.prop")
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(prop)
	scanner.Split(bufio.ScanLines)
	matchedVer := false
	matchedTarget := false
	var version vars.Version
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, "ro.anki.victor.version") {
			matchedVer = true
			ver := strings.Split(text, "=")[1]
			version, err = manager.SplitVersion(ver)
			if err != nil {
				return err
			}
		}
		if strings.Contains(text, "ro.build.target") {
			matchedTarget = true
			target, err = strconv.Atoi(strings.Split(text, "=")[1])
			if err != nil {
				return err
			}
		}
	}
	if !matchedVer {
		return errors.New("version not in build.prop")
	}
	if !matchedTarget {
		return errors.New("target not found in build.prop, you must run patcher")
	}
	fmt.Println("Found version " + version.Full + " with target " + vars.Targets[target])
	prop.Close()
	fmt.Println("Unmounting image...")
	err = UnmountImage("./work")
	if err != nil {
		return err
	}
	fmt.Println("Compressing sysfs...")
	imageBytes, _ := os.ReadFile("./tmp/apq8009-robot-sysfs.img")
	compressed, err := CompressBytes(imageBytes)
	if err != nil {
		return err
	}
	fmt.Println("Encrypting sysfs...")
	encrypted, err := EncryptBytes(compressed)
	if err != nil {
		return err
	}
	fmt.Println("Packing tar...")
	err = PackTar(encrypted, target, CreateManifest(version, "./tmp/apq8009-robot-sysfs.img", target), version)
	if err != nil {
		return err
	}
	fmt.Println("Success packing OTA!")
	return nil
}
