package install

import (
	"fmt"
	"os"

	"appimage-manager/app/commands"
	"appimage-manager/app/repos"
	"appimage-manager/app/utils"
)

type InstallCmd struct {
	Target string `arg name:"target" help:"Installation target." type:"string"`
}

func (cmd *InstallCmd) Run(*commands.Context) (err error) {
	repo, err := repos.ParseTarget(cmd.Target)
	if err != nil {
		return err
	}

	release, err := repo.GetLatestRelease()
	if err != nil {
		return err
	}

	selectedBinary, err := utils.PromptBinarySelection(release.Files)
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

	cmd.createDesktopIntegration(err, targetFilePath)

	return
}

func (cmd *InstallCmd) addToRegistry(targetFilePath string, repo repos.Repo) {
	sha1, _ := utils.GetFileSHA1(targetFilePath)
	updateInfo, _ := utils.ReadUpdateInfo(targetFilePath)
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

func (cmd *InstallCmd) createDesktopIntegration(err error, targetFilePath string) {
	libAppImage, err := utils.NewLibAppImageBindings()
	if err != nil {
		fmt.Println("Integration failed: missing libappimage.so")
		return
	}

	fmt.Println("Creating menu entry and mime-types integrations")
	err = libAppImage.Register(targetFilePath)
	if err != nil {
		fmt.Println("Integration failed: " + err.Error())
	} else {
		fmt.Println("Integration completed")
	}
}
