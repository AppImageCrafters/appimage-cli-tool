package commands

// #include <stdio.h>
// #include <stdlib.h>
import "C"

import (
	"appimage-installer/app/utils"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/manifoldco/promptui"
	"github.com/rainycape/dl"
)

type InstallCmd struct {
	Id string `arg name:"id" help:"Installation id." type:"string"`
}

type DownloadLink struct {
	Name string
	Url  string
}

func (cmd *InstallCmd) Run(*Context) (err error) {
	idParts := strings.SplitN(cmd.Id, ":", 2)
	if len(idParts) < 2 {
		return fmt.Errorf("invalid installation id '%s'", cmd.Id)
	}

	source := idParts[0]

	switch source {
	case "appimagehub":
		err = cmd.appImageHubInstall(idParts[1])
	default:
		return fmt.Errorf("invalid installation id '%s'", cmd.Id)
	}

	return
}

func (cmd *InstallCmd) appImageHubInstall(path string) (err error) {
	doc, err := xmlquery.LoadURL("https://www.appimagehub.com/ocs/v1/content/data/" + path)
	if err != nil {
		return
	}

	downloadLinks, err := cmd.appImageHubParseDownloadLinks(doc)
	if err != nil {
		return err
	}

	result, err := cmd.promptBinarySelection(downloadLinks)
	if err != nil {
		return err
	}

	filePath, err := cmd.makeTargetFilePath(result)
	if err != nil {
		return err
	}

	err = utils.DownloadAppImage(result.Url, filePath)
	if err != nil {
		return err
	}
	fmt.Println("AppImage downloaded to: " + filePath)

	registry, _ := utils.OpenRegistry()
	if registry != nil {
		_ = registry.Set(result.Name, "appimagehub:"+path)
		_ = registry.Close()
	}

	err = installAppImage(filePath)
	if err != nil {
		fmt.Println("Registration failed: " + err.Error())
	} else {
		fmt.Println("Registration completed")
	}
	return
}

func (cmd *InstallCmd) appImageHubParseDownloadLinks(doc *xmlquery.Node) ([]DownloadLink, error) {
	var downloadLinks []DownloadLink
	for i := 1; i < 100; i++ {
		idx := strconv.Itoa(i)
		link, err := xmlquery.Query(doc, "//ocs/data/content/downloadlink"+idx+"/text()")
		if err != nil {
			return nil, err
		}
		name, err := xmlquery.Query(doc, "//ocs/data/content/downloadname"+idx+"/text()")
		if err != nil {
			return nil, err
		}

		if link == nil {
			break
		}

		downloadLink := DownloadLink{
			Name: name.Data,
			Url:  link.Data,
		}

		if strings.HasSuffix(downloadLink.Name, ".AppImage") {
			downloadLinks = append(downloadLinks, downloadLink)
		}
	}
	return downloadLinks, nil
}

func (cmd *InstallCmd) promptBinarySelection(downloadLinks []DownloadLink) (result *DownloadLink, err error) {
	if len(downloadLinks) == 1 {
		return &downloadLinks[0], nil
	}

	prompt := promptui.Select{
		Label: "Select binary to download",
		Items: downloadLinks,
		Templates: &promptui.SelectTemplates{
			Label:    "   {{ .Name }}",
			Active:   "\U00002705 {{ .Name }}",
			Inactive: "   {{ .Name }}",
			Selected: "\U00002705 {{ .Name }}"},
	}

	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return &downloadLinks[i], nil
}

func installAppImage(filePath string) error {
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

	var appimage_register_in_system func(path *C.char, verbose int) int
	err = lib.Sym("appimage_register_in_system", &appimage_register_in_system)
	if err != nil {
		return err
	}

	if appimage_shall_not_be_integrated(C.CString(filePath)) != 0 {
		return nil
	}

	if appimage_register_in_system(C.CString(filePath), 1) != 0 {
		return fmt.Errorf("registration failed")
	}

	return nil
}

func (cmd *InstallCmd) makeTargetFilePath(link *DownloadLink) (string, error) {
	applicationsPath, err := utils.MakeApplicationsDirPath()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(applicationsPath, link.Name)
	return filePath, nil
}
