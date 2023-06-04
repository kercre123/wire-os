package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kercre123/wire-os/wired/pkg/modify"
	"github.com/kercre123/wire-os/wired/pkg/vars"
)

func main() {
	// set modifiers to wireos
	vars.Modifiers = modify.WireOSModifiers

	fmt.Println("wire_d running, running init functions")
	vars.Init()
	fmt.Println("starting shell (debug)")
	for {
		fmt.Printf("\nEnter a command: ")
		var command string
		fmt.Scanln(&command)
		switch {
		case strings.Contains(command, "list"):
			fmt.Println(strings.TrimSpace(string(modify.GetAllModifiers())))
		case strings.Contains(command, "apply"):
			fmt.Println("")
			fmt.Println("All modifiers:")
			for ind, mod := range vars.Modifiers {
				fmt.Println(fmt.Sprint(ind) + ": " + mod.Name + " (" + mod.Description + ")")
			}
			fmt.Println("For reference, modifiers in loaded modifiers file:")
			for _, mod := range vars.LoadedModifiers {
				fmt.Println(fmt.Sprint(mod.ModifierID) + ": " + vars.Modifiers[mod.ModifierID].Name + " (Applied: " + fmt.Sprint(mod.Applied) + ")")
			}
			fmt.Printf("\nWhich modifier would you like to apply? (ex. 0): ")
			var modApply string
			fmt.Scanln(&modApply)
			modID, err := strconv.Atoi(modApply)
			if err != nil {
				fmt.Println("Given value was not an int")
				continue
			}
			err = modify.ApplyModifier(modID)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case strings.Contains(command, "remove"):
			fmt.Println("")
			fmt.Println("All modifiers:")
			for ind, mod := range vars.Modifiers {
				fmt.Println(fmt.Sprint(ind) + ": " + mod.Name + " (" + mod.Description + ")")
			}
			fmt.Println("For reference, modifiers in loaded modifiers file:")
			for _, mod := range vars.LoadedModifiers {
				fmt.Println(fmt.Sprint(mod.ModifierID) + ": " + vars.Modifiers[mod.ModifierID].Name + " (Applied: " + fmt.Sprint(mod.Applied) + ")")
			}
			fmt.Printf("\nWhich modifier would you like to remove? (ex. 0): ")
			var modRemove string
			fmt.Scanln(&modRemove)
			modID, err := strconv.Atoi(modRemove)
			if err != nil {
				fmt.Println("Given value was not an int")
				continue
			}
			err = modify.RemoveModifier(modID)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}
