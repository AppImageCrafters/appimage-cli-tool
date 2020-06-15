package install

import (
	"fmt"
	"strconv"
	"strings"

	"appimage-manager/app/commands"
	"appimage-manager/app/utils"

	"github.com/antchfx/xmlquery"
)

type InstallCmd struct {
	Target string `arg name:"target" help:"Installation target." type:"string"`
}

func (cmd *InstallCmd) Run(*commands.Context) (err error) {
	cmd.Target, err = utils.UrlToTarget(cmd.Target)
	if err != nil {
		return err
	}

	targetParts := strings.SplitN(cmd.Target, ":", 2)
	if len(targetParts) < 2 {
		return fmt.Errorf("invalid installation id '%s'", cmd.Target)
	}

	source := targetParts[0]

	switch source {
	case "appimagehub":
		err = cmd.appImageHubInstall(targetParts[1])
	case "github":
		err = InstallGithubTarget(targetParts[1])
	default:
		return fmt.Errorf("invalid installation id '%s'", cmd.Target)
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

	result, err := utils.PromptBinarySelection(downloadLinks)
	if err != nil {
		return err
	}

	filePath, err := utils.MakeTargetFilePath(result)
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

	err = utils.InstallAppImage(filePath)
	if err != nil {
		fmt.Println("Registration failed: " + err.Error())
	} else {
		fmt.Println("Registration completed")
	}
	return
}

func (cmd *InstallCmd) appImageHubParseDownloadLinks(doc *xmlquery.Node) ([]utils.DownloadLink, error) {
	var downloadLinks []utils.DownloadLink
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

		downloadLink := utils.DownloadLink{
			Name: name.Data,
			Url:  link.Data,
		}

		if strings.HasSuffix(downloadLink.Name, ".AppImage") {
			downloadLinks = append(downloadLinks, downloadLink)
		}
	}
	return downloadLinks, nil
}
