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
	Target string `arg name:"target" help:"Updates the target application." type:"string"`
}

func (cmd *UpdateCmd) Run(*Context) (err error) {
	filePath, err := cmd.getBundleFilePath()
	if err != nil {
		return err
	}

	updateMethod, err := update.NewUpdaterFor(filePath)
	if err != nil {
		return err
	}

	fmt.Println("Looking for updates of: ", filePath)
	updateAvailable, err := updateMethod.Lookup()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !updateAvailable {
		fmt.Println("No updates were found for: ", filePath)
		return
	}

	result, err := updateMethod.Download()
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}

	fmt.Println("Update downloaded to: " + result)

	return nil
}

func (cmd *UpdateCmd) getBundleFilePath() (string, error) {
	if strings.HasPrefix(cmd.Target, "file://") {
		cmd.Target = cmd.Target[7:]
	}

	if _, err := os.Stat(cmd.Target); err == nil {
		return cmd.Target, nil
	}

	registry, err := utils.OpenRegistry()
	if err != nil {
		return "", err
	}
	registry.Update()

	fileName, ok := registry.Lookup(cmd.Target)
	if !ok {
		fileName = cmd.Target
	}

	applicationsDir, err := utils.MakeApplicationsDirPath()
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(applicationsDir, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("application not found \"" + cmd.Target + "\"")
	}
	return filePath, nil
}
