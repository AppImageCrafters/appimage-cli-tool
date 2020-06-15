package install

import (
	"context"
	"fmt"
	"strings"

	"appimage-manager/app/utils"

	"github.com/google/go-github/v31/github"
)

func InstallGithubTarget(target string) error {
	target_parts := strings.Split(target, "/")

	download_links, err := getLatestReleaseDownloadLinks(target_parts)
	if err != nil {
		return err
	}

	if len(download_links) == 0 {
		return fmt.Errorf("no AppImages releases found")
	}

	download_link, err := utils.PromptBinarySelection(download_links)

	target_path, err := utils.MakeTargetFilePath(download_link)
	if err != nil {
		return err
	}

	err = utils.DownloadAppImage(download_link.Url, target_path)
	if err != nil {
		return err
	}

	registry, _ := utils.OpenRegistry()
	if registry != nil {
		_ = registry.Set(download_link.Name, "github:"+target)
		_ = registry.Close()
	}

	err = utils.InstallAppImage(target_path)
	if err != nil {
		fmt.Println("Registration failed: " + err.Error())
	} else {
		fmt.Println("Registration completed")
	}

	return nil
}
func getLatestReleaseDownloadLinks(target_parts []string) ([]utils.DownloadLink, error) {
	download_links := []utils.DownloadLink{}

	client := github.NewClient(nil)
	releases, _, err := client.Repositories.ListReleases(context.Background(), target_parts[0], target_parts[1], nil)
	if err != nil {
		return nil, err
	}

	for _, release := range releases {
		if *release.Draft == true {
			continue
		}

		for _, asset := range release.Assets {
			if strings.HasSuffix(*asset.Name, ".AppImage") {
				download_links = append(download_links, utils.DownloadLink{
					Name: *asset.Name,
					Url:  *asset.BrowserDownloadURL,
				})
			}
		}

		if len(download_links) > 0 {
			return download_links, nil
		}
	}

	return download_links, nil
}
