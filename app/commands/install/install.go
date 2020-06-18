package install

import (
	"appimage-manager/app/commands"
	"appimage-manager/app/utils"
	"fmt"
	updateUtils "github.com/AppImageCrafters/appimage-update/util"
	"os"
)

type InstallCmd struct {
	Target string `arg name:"target" help:"Installation target." type:"string"`
}

type Release struct {
	Tag   string
	files []utils.BinaryUrl
}

type Repo interface {
	Id() string
	GetLatestRelease() (*Release, error)
	Download(binaryUrl *utils.BinaryUrl, targetPath string) error
	FallBackUpdateInfo() string
}

func (cmd *InstallCmd) Run(*commands.Context) (err error) {
	repo, err := cmd.parseTarget()
	if err != nil {
		return err
	}

	release, err := repo.GetLatestRelease()
	if err != nil {
		return err
	}

	selectedBinary, err := utils.PromptBinarySelection(release.files)
	if err != nil {
		return err
	}

	targetFilePath, err := utils.MakeTargetFilePath(selectedBinary)
	if err != nil {
		return err
	}

	if _, err = os.Stat(targetFilePath); err == nil {
		return ApplicationInstalled
	}

	err = repo.Download(selectedBinary, targetFilePath)
	if err != nil {
		return err
	}

	cmd.addToRegistry(targetFilePath, repo)

	err = cmd.makeDesktopIntegration(err, targetFilePath)
	return
}

func (cmd *InstallCmd) parseTarget() (Repo, error) {
	repo, err := NewGitHubRepo(cmd.Target)
	if err == nil {
		return repo, nil
	}

	repo, err = NewAppImageHubRepo(cmd.Target)
	if err == nil {
		return repo, nil
	}

	return nil, InvalidTargetFormat
}

func (cmd *InstallCmd) addToRegistry(targetFilePath string, repo Repo) {
	sha1, _ := utils.GetFileSHA1(targetFilePath)
	updateInfo, _ := updateUtils.ReadUpdateInfo(targetFilePath)
	if updateInfo == "" {
		updateInfo = repo.FallBackUpdateInfo()
	}

	entry := utils.RegistryEntry{
		FilePath:   targetFilePath,
		Repo:       repo.Id(),
		FileSha1:   sha1,
		UpdateInfo: updateInfo,
	}

	registry, _ := utils.OpenRegistry()
	if registry != nil {
		_ = registry.Add(entry)
		_ = registry.Close()
	}
}

func (cmd *InstallCmd) makeDesktopIntegration(err error, targetFilePath string) error {
	fmt.Println("Integrating with the desktop environment")
	err = utils.Integrate(targetFilePath)
	if err != nil {
		fmt.Println("Integration failed: " + err.Error())
	} else {
		fmt.Println("Integration completed")
	}
	return err
}
