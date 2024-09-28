package handler

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/kercre123/wire-os/pkg/vars"
)

const (
	PassphraseFile = "./resources/certs/ota.pas"
)

func DecryptBytes(in []byte) ([]byte, error) {
	inFileName := os.TempDir() + "/wireos-encrypted.gz"
	outFileName := os.TempDir() + "/wireos-decrypted.gz"
	// i cannot for the life of me figure out how to do this in pure Go
	// openssl enc -d -aes-256-ctr -pass file:certs/ota.pas -in test.gz -out testout.gz
	err := os.WriteFile(inFileName, in, 0777)
	if err != nil {
		return in, err
	}
	err = exec.Command("openssl", "enc", "-d", "-aes-256-ctr", "-pass", "file:"+PassphraseFile, "-md", "md5", "-in", inFileName, "-out", outFileName).Run()
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

func DecompressBytes(in []byte) ([]byte, error) {
	reader := bytes.NewReader(in)
	gzreader, err := gzip.NewReader(reader)
	if err != nil {
		return []byte{}, err
	}
	decompressed, err := io.ReadAll(gzreader)
	if err != nil {
		return []byte{}, err
	}
	return decompressed, nil
}

func GetOTASystem(filename string) ([]byte, error) {
	ota, err := os.Open(filename)
	if err != nil {
		return []byte{}, err
	}
	reader := tar.NewReader(ota)
	for {
		head, err := reader.Next()
		if err == io.EOF {
			return []byte{}, errors.New("no system gz in tar")
		} else if err != nil {
			return []byte{}, err
		}
		if strings.Contains(head.Name, "sysfs") {
			return io.ReadAll(reader)
		}
	}
}

func MountImage(image string, path string) error {
	//guestmount -a thing.img -m /dev/sda1 --rw edits/
	//out, err := exec.Command("guestmount", "-a", image, "-m", "/dev/sda", "--rw", "-o", "nonempty", "-o", "direct_io", "-o", "noatime", path).Output()
	// if err != nil {
	// 	fmt.Println(string(out))
	// 	return err
	// }
	err := exec.Command("mount", "-o", "rw,nodiratime,noatime", image, path).Run()
	if err != nil {
		return err
	}
	return nil
}

func MountOTA(version vars.Version) error {
	image := "./store/" + version.Full + ".ota"
	system, err := GetOTASystem(image)
	if err != nil {
		return err
	}
	fmt.Println("Found OTA system GZ, decrypting...")
	decrypted, err := DecryptBytes(system)
	if err != nil {
		return err
	}
	fmt.Println("Successfully decrypted, decompressing...")
	decompressed, err := DecompressBytes(decrypted)
	if err != nil {
		return err
	}
	fmt.Println("Decompressed, writing file...")
	err = os.WriteFile("./tmp/apq8009-robot-sysfs.img", decompressed, 0777)
	if err != nil {
		return err
	}
	fmt.Println("File written, mounting...")
	err = MountImage("./tmp/apq8009-robot-sysfs.img", "./work")
	if err != nil {
		return err
	}
	fmt.Println("Successfully mounted to ./work!")
	return nil
}
