package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"appimage-installer/app/utils"

	"github.com/AppImageCrafters/appimage-update"
)

type UpdateCmd struct {
	Id string `arg name:"id" help:"Installation id or file name." type:"string"`
}

func (cmd *UpdateCmd) Run(*Context) (err error) {
	filePath, err := cmd.getApplicationFilePath()
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

func (cmd *UpdateCmd) getApplicationFilePath() (string, error) {
	registry, err := utils.OpenRegistry()
	if err != nil {
		return "", err
	}
	registry.Update()

	fileName, ok := registry.Lookup(cmd.Id)
	if !ok {
		fileName = cmd.Id
	}

	applicationsDir, err := utils.MakeApplicationsDirPath()
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(applicationsDir, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("application not found \"" + cmd.Id + "\"")
	}
	return filePath, nil
}
