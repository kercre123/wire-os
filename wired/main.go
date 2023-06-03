package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	fmt.Println("Hello world!")
	out, err := exec.Command("uname", "-a").Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("uname: " + strings.TrimSpace(string(out)))
}
