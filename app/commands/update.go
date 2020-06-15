package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"appimage-manager/app/utils"

	"github.com/AppImageCrafters/appimage-update"
)

type UpdateCmd struct {
	Targets []string `arg optional name:"targets" help:"Updates the target applications." type:"string"`

	Check bool `help:"Only check for updates."`
	All   bool `help:"Update all applications."`
}

func (cmd *UpdateCmd) Run(*Context) (err error) {
	if cmd.All {
		cmd.Targets, err = getAllTargets()
		if err != nil {
			return err
		}
	}

	for _, target := range cmd.Targets {
		filePath, err := cmd.getBundleFilePath(target)
		if err != nil {
			println(err.Error())
			continue
		}

		updateMethod, err := update.NewUpdaterFor(filePath)
		if err != nil {
			println(err.Error())
			continue
		}

		fmt.Println("Looking for updates of: ", filePath)
		updateAvailable, err := updateMethod.Lookup()
		if err != nil {
			println(err.Error())
			continue
		}

		if !updateAvailable {
			fmt.Println("No updates were found for: ", filePath)
			continue
		}

		if cmd.Check {
			fmt.Println("Update available for: ", filePath)
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

func getAllTargets() ([]string, error) {
	registry, err := utils.OpenRegistry()
	if err != nil {
		return nil, err
	}
	registry.Update()

	paths := make([]string, len(registry.Entries))
	for k := range registry.Entries {
		paths = append(paths, k)
	}

	return paths, nil
}

func (cmd *UpdateCmd) getBundleFilePath(target string) (string, error) {
	if strings.HasPrefix(target, "file://") {
		cmd.Targets = cmd.Targets[7:]
	}

	if _, err := os.Stat(target); err == nil {
		return target, nil
	}

	registry, err := utils.OpenRegistry()
	if err != nil {
		return "", err
	}
	registry.Update()

	fileName, ok := registry.Lookup(target)
	if !ok {
		fileName = target
	}

	applicationsDir, err := utils.MakeApplicationsDirPath()
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(applicationsDir, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("application not found \"" + target + "\"")
	}
	return filePath, nil
}
