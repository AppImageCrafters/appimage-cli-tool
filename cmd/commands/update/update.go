package update

import (
	"fmt"
	"strings"

	"github.com/AppImageCrafters/appimage-update"
	"github.com/AppImageCrafters/appimage-update/updaters"
)

func NewUpdater(updateInfoString string, appImagePath string) (update.Updater, error) {
	if strings.HasPrefix(updateInfoString, "zsync") {
		return updaters.NewZSyncUpdater(&updateInfoString, appImagePath)
	}

	if strings.HasPrefix(updateInfoString, "gh-releases-zsync") {
		return updaters.NewGitHubZsyncUpdater(&updateInfoString, appImagePath)
	}

	if strings.HasPrefix(updateInfoString, "gh-releases-direct") {
		return updaters.NewGitHubDirectUpdater(&updateInfoString, appImagePath)
	}

	if strings.HasPrefix(updateInfoString, "ocs-v1-appimagehub-direct") {
		return updaters.NewOCSAppImageHubDirect(&updateInfoString, appImagePath)
	}

	if strings.HasPrefix(updateInfoString, "ocs-v1-appimagehub-zsync") {
		return updaters.NewOCSAppImageHubZSync(&updateInfoString, appImagePath)
	}

	return nil, fmt.Errorf("Invalid updated information: ", updateInfoString)
}
