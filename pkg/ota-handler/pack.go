package handler

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"os/exec"

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
	err := exec.Command("umount", path).Run()
	if err != nil {
		return err
	}
	return nil
}

func PackTar(systembytes []byte, target int, manifest []byte, version vars.Version) error {
	// 0 = dev, 1 = whiskey, 2 = oskr, 3 = orange(fills a very specific niche of bot)
	var buf []byte
	if target > 3 || target < 0 {
		return errors.New("target out of range, must be between 1 and 3")
	}
	bootbytes, err := os.ReadFile("./boots/" + vars.Targets[target] + ".img.gz")
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
	if err != nil {
		return err
	}
	// sha would go here
	err = tarFile.WriteHeader(&tar.Header{
		Name: "apq8009-robot-boot.img.gz",
		Size: int64(len(bootbytes)),
		Mode: 0777,
	})
	if err != nil {
		return err
	}
	tarFile.Write(systembytes)
	if err != nil {
		return err
	}
	err = tarFile.WriteHeader(&tar.Header{
		Name: "apq8009-robot-sysfs.img.gz",
		Size: int64(len(systembytes)),
		Mode: 0777,
	})
	if err != nil {
		return err
	}
	tarFile.Write(systembytes)
	if err != nil {
		return err
	}
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

func PackOTA(version vars.Version, target int) error {
	if target > 3 || target < 0 {
		return errors.New("target out of range, must be between 1 and 3")
	}
	fmt.Println("Unmounting image...")
	err := UnmountImage("./work")
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
