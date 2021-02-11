package commands

import (
	"appimage-cli-tool/internal/utils"
	"fmt"
	"os"
)

type RemoveCmd struct {
	Target   string `arg name:"id" help:"Installation id or file name." type:"string"`
	KeepFile bool   `help:"Remove only the application desktop entry."`
}

func (cmd *RemoveCmd) Run(*Context) (err error) {
	registry, err := utils.OpenRegistry()
	if err != nil {
		return err
	}
	defer registry.Close()

	registry.Update()

	entry, ok := registry.Lookup(cmd.Target)
	if !ok {
		return fmt.Errorf("application not found \"" + cmd.Target + "\"")
	}

	err = removeDesktopIntegration(entry.FilePath)
	if err != nil {
		fmt.Println("Desktop integration removal failed: " + err.Error())
	}

	if cmd.KeepFile {
		return nil
	}

	err = os.Remove(entry.FilePath)
	if err != nil {
		return fmt.Errorf("Unable to remove AppImage file: " + err.Error())
	}
	fmt.Println("Application removed: " + entry.FilePath)
	registry.Remove(entry.FilePath)
	return err
}

func removeDesktopIntegration(filePath string) error {
	libAppImage, err := utils.NewLibAppImageBindings()
	if err != nil {
		return err
	}

	if libAppImage.ShallNotBeIntegrated(filePath) {
		return nil
	}

	return libAppImage.Unregister(filePath)
}
