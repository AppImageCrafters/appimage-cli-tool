package repos

import (
	"appimage-cli-tool/internal/utils"
)

type Release struct {
	Tag   string
	Files []utils.BinaryUrl
}

type Repo interface {
	Id() string
	GetLatestRelease() (*Release, error)
	Download(binaryUrl *utils.BinaryUrl, targetPath string) error
	FallBackUpdateInfo() string
}

func ParseTarget(target string) (Repo, error) {
	repo, err := NewGitHubRepo(target)
	if err == nil {
		return repo, nil
	}

	repo, err = NewAppImageHubRepo(target)
	if err == nil {
		return repo, nil
	}

	return nil, InvalidTargetFormat
}
