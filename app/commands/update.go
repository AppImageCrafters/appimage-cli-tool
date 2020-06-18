package commands

import (
	"appimage-manager/app/utils"
	"errors"
	"fmt"
	updateUtils "github.com/AppImageCrafters/appimage-update/util"

	"github.com/AppImageCrafters/appimage-update"
)

type UpdateCmd struct {
	Targets []string `arg optional name:"targets" help:"Updates the target applications." type:"string"`

	Check bool `help:"Only check for updates."`
	All   bool `help:"Update all applications."`
}

var NoUpdateInfo = errors.New("there is no update information")

func (cmd *UpdateCmd) Run(*Context) (err error) {
	if cmd.All {
		cmd.Targets, err = getAllTargets()
		if err != nil {
			return err
		}
	}

	for _, target := range cmd.Targets {
		entry, err := cmd.getRegistryEntry(target)
		if err != nil {
			continue
		}

		updateMethod, err := update.NewUpdateForUpdateString(entry.UpdateInfo, entry.FilePath)
		if err != nil {
			println(err.Error())
			continue
		}

		fmt.Println("Looking for updates of: ", entry.FilePath)
		updateAvailable, err := updateMethod.Lookup()
		if err != nil {
			println(err.Error())
			continue
		}

		if !updateAvailable {
			fmt.Println("No updates were found for: ", entry.FilePath)
			continue
		}

		if cmd.Check {
			fmt.Println("Update available for: ", entry.FilePath)
			continue
		}

		result, err := updateMethod.Download()
		if err != nil {
			println(err.Error())
			continue
		}

		fmt.Println("Update downloaded to: " + result)
	}

	return nil
}

func (cmd *UpdateCmd) getRegistryEntry(target string) (utils.RegistryEntry, error) {
	registry, err := utils.OpenRegistry()
	if err != nil {
		return utils.RegistryEntry{}, err
	}
	defer registry.Close()

	entry, _ := registry.Lookup(target)

	if entry.UpdateInfo == "" {
		entry.UpdateInfo, _ = updateUtils.ReadUpdateInfo(target)
		entry.FilePath = target
	}

	if entry.UpdateInfo == "" {
		return entry, NoUpdateInfo
	} else {
		return entry, nil
	}
}

func getAllTargets() ([]string, error) {
	registry, err := utils.OpenRegistry()
	if err != nil {
		return nil, err
	}
	registry.Update()

	var paths []string
	for k := range registry.Entries {
		paths = append(paths, k)
	}

	return paths, nil
}
