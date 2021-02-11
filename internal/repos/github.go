package repos

import (
	"appimage-cli-tool/internal/utils"
	"context"
	"strings"
	"github.com/google/go-github/v31/github"
)

type GitHubRepo struct {
	User    string
	Project string
	Release string
	File    string
}

func NewGitHubRepo(target string) (repo Repo, err error) {
	if strings.HasPrefix(target, "https://github.com/") {
		target = strings.Replace(target, "https://github.com/", "github:", 1)
	}

	if !strings.HasPrefix(target, "github:") {
		return repo, InvalidTargetFormat
	}

	repo = &GitHubRepo{}

	targetParts := strings.Split(target[7:], "/")
	targetPartsLen := len(targetParts)
	if targetPartsLen < 2 {
		return repo, InvalidTargetFormat
	}

	ghSource := GitHubRepo{
		User:    targetParts[0],
		Project: targetParts[1],
	}

	if targetPartsLen > 2 {
		ghSource.Release = targetParts[2]
	}

	if targetPartsLen > 3 {
		ghSource.File = targetParts[3]
	}

	return &ghSource, nil
}

func (g GitHubRepo) Id() string {
	id := "github:" + g.User + "/" + g.Project

	return id
}

func (g GitHubRepo) GetLatestRelease() (*Release, error) {
	var downloadLinks []utils.BinaryUrl

	client := github.NewClient(nil)
	releases, _, err := client.Repositories.ListReleases(context.Background(), g.User, g.Project, nil)
	if err != nil {
		return nil, err
	}

	for _, release := range releases {
		if *release.Draft == true {
			continue
		}

		for _, asset := range release.Assets {
			if strings.HasSuffix(*asset.Name, ".AppImage") {
				downloadLinks = append(downloadLinks, utils.BinaryUrl{
					FileName: *asset.Name,
					Url:      *asset.BrowserDownloadURL,
				})
			}
		}

		if len(downloadLinks) > 0 {
			return &Release{
				*release.TagName,
				downloadLinks,
			}, nil
		}
	}

	return nil, NoAppImageBinariesFound
}

func (g GitHubRepo) Download(binaryUrl *utils.BinaryUrl, targetPath string) (err error) {
	err = utils.DownloadAppImage(binaryUrl.Url, targetPath)
	return
}

func (g GitHubRepo) FallBackUpdateInfo() string {
	updateInfo := "gh-releases-direct|" + g.User + "|" + g.Project
	if g.Release == "" {
		updateInfo += "|latest"
	} else {
		updateInfo += "|" + g.Release
	}

	if g.File == "" {
		updateInfo += "|*.AppImage"
	} else {
		updateInfo += "|" + g.File
	}

	return updateInfo
}
