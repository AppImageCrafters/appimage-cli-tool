package commands

// #include <stdio.h>
// #include <stdlib.h>
import "C"
import (
	"appimage-installer/app/utils"
	"fmt"
	"github.com/rainycape/dl"
	"os"
	"path/filepath"
)

type RemoveCmd struct {
	Id string `arg name:"id" help:"Installation id or file name." type:"string"`
}

func (cmd *RemoveCmd) Run(*Context) (err error) {
	registry, err := utils.OpenRegistry()
	if err != nil {
		return err
	}
	registry.Update()

	fileName, ok := cmd.lookupFileName(registry)
	if !ok {
		fileName = cmd.Id
	}

	applicationsDir, err := utils.MakeApplicationsDirPath()
	if err != nil {
		return err
	}
	filePath := filepath.Join(applicationsDir, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("application not found \"" + cmd.Id + "\"")
	}

	err = uninstallAppImage(filePath)
	if err != nil {
		fmt.Println("Desktop deregistration failed: " + err.Error())
	}

	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("Unable to remove AppImage file: " + err.Error())
	}
	fmt.Println("Application removed: " + fileName)
	return err
}

func (cmd *RemoveCmd) lookupFileName(registry *utils.Registry) (string, bool) {
	for fileName, entry := range registry.Entries {
		if entry.Id == cmd.Id {
			return fileName, true
		}
	}

	return "", false
}

func uninstallAppImage(filePath string) error {
	lib, err := dl.Open("libappimage.so", 0)
	if err != nil {
		return fmt.Errorf("desktop integration not available")
	}
	defer lib.Close()

	var appimage_shall_not_be_integrated func(path *C.char) int
	err = lib.Sym("appimage_shall_not_be_integrated", &appimage_shall_not_be_integrated)
	if err != nil {
		return err
	}

	var appimage_unregister_in_system func(path *C.char, verbose int) int
	err = lib.Sym("appimage_unregister_in_system", &appimage_unregister_in_system)
	if err != nil {
		return err
	}

	if appimage_shall_not_be_integrated(C.CString(filePath)) != 0 {
		return nil
	}

	if appimage_unregister_in_system(C.CString(filePath), 1) != 0 {
		return fmt.Errorf("deregistration failed")
	}

	return nil
}
