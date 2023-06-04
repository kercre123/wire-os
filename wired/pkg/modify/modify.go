package modify

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kercre123/wire-os/wired/pkg/vars"
)

// Returns all modifiers as JSON []byte
func GetAllModifiers() []byte {
	var repMods []vars.ReportedModifier
	for ind, mod := range vars.Modifiers {
		repMods = append(repMods, vars.ReportedModifier{
			Name:        mod.Name,
			Description: mod.Description,
			ModifierID:  ind,
			Applied:     false,
		})
	}
	for _, mod := range vars.LoadedModifiers {
		if mod.Applied {
			repMods[mod.ModifierID].Applied = true
		}
	}
	jsonBytes, err := json.Marshal(repMods)
	if err != nil {
		fmt.Println("Error marshalling reported modifiers: " + err.Error())
		return []byte(err.Error())
	}
	return jsonBytes
}

func GetLoadedModifier(ID int) (vars.LoadedModifier, error) {
	for _, lmod := range vars.LoadedModifiers {
		if lmod.ModifierID == ID {
			return lmod, nil
		}
	}
	return vars.LoadedModifier{
		Applied: false,
	}, errors.New("loaded modifier does not exist")
}

func GetModifier(ID int) (vars.Modifier, error) {
	if ID > len(vars.Modifiers)-1 {
		return vars.Modifier{}, errors.New("id out of range")
	}
	return vars.Modifiers[ID], nil
}

func CreateLoadedModifier(ModifierID int) {
	// false until after application. this is run before
	lmod := vars.LoadedModifier{
		ModifierID: ModifierID,
		Applied:    false,
	}
	vars.LoadedModifiers = append(vars.LoadedModifiers, lmod)
	vars.SaveModifiersJSON()
}

func SetLoadedModifier(ModifierID int, applied bool) error {
	for ind, lmod := range vars.LoadedModifiers {
		if lmod.ModifierID == ModifierID {
			vars.LoadedModifiers[ind].Applied = applied
			vars.SaveModifiersJSON()
			return nil
		}
	}
	return errors.New("setloadedmodifier did not match")
}

func ApplyModifier(ID int) error {
	modifier, err := GetModifier(ID)
	if err != nil {
		return err
	}
	lmod, err := GetLoadedModifier(ID)
	if err != nil {
		CreateLoadedModifier(ID)
	}
	if lmod.Applied {
		return errors.New("modifier already applied")
	}
	fmt.Println("Applying modifier " + modifier.Name + " (" + modifier.Description + ")")
	err = modifier.Apply()
	if err != nil {
		return err
	}
	SetLoadedModifier(ID, true)
	fmt.Println(modifier.Name + " applied successfully!")
	return nil
}

func RemoveModifier(ID int) error {
	modifier, err := GetModifier(ID)
	if err != nil {
		return err
	}
	lmod, err := GetLoadedModifier(ID)
	if err != nil {
		return errors.New("modifier must be applied first")
	}
	if !lmod.Applied {
		return errors.New("modifier already removed")
	}
	fmt.Println("Removing modifier " + modifier.Name + " (" + modifier.Description + ")")
	err = modifier.Remove()
	if err != nil {
		return err
	}
	SetLoadedModifier(ID, false)
	fmt.Println(modifier.Name + " removed successfully!")
	return nil
}
