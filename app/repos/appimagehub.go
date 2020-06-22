package repos

import (
	"appimage-manager/app/utils"
	"github.com/antchfx/xmlquery"
	"strconv"
	"strings"
)

type AppImageHubRepo struct {
	ContentId string
}

func NewAppImageHubRepo(target string) (Repo, error) {
	if strings.HasPrefix(target, "https://www.appimagehub.com/p/") {
		target = strings.Replace(target, "https://www.appimagehub.com/p/", "appimagehub:", 1)
	}

	if !strings.HasPrefix(target, "appimagehub:") {
		return nil, InvalidTargetFormat
	}

	repo := &AppImageHubRepo{}
	repo.ContentId = target[12:]

	return repo, nil
}

func (a AppImageHubRepo) Id() string {
	return "appimagehub:" + a.ContentId
}

func (a AppImageHubRepo) GetLatestRelease() (*Release, error) {
	doc, err := xmlquery.LoadURL("https://www.appimagehub.com/ocs/v1/content/data/" + a.ContentId)
	if err != nil {
		return nil, err
	}

	var downloadLinks []utils.BinaryUrl
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

		downloadLink := utils.BinaryUrl{
			FileName: name.Data,
			Url:      link.Data,
		}

		if strings.HasSuffix(downloadLink.FileName, ".AppImage") ||
			strings.HasSuffix(downloadLink.FileName, ".appimage") {
			downloadLinks = append(downloadLinks, downloadLink)
		}
	}

	if len(downloadLinks) > 0 {
		return &Release{
			"latest",
			downloadLinks,
		}, nil
	} else {
		return nil, NoAppImageBinariesFound
	}
}

func (a AppImageHubRepo) Download(binaryUrl *utils.BinaryUrl, targetPath string) (err error) {
	err = utils.DownloadAppImage(binaryUrl.Url, targetPath)
	return
}

func (a AppImageHubRepo) FallBackUpdateInfo() string {
	return "ocs-v1-appimagehub-zsync|www.appimagehub.com/ocs/v1|" + a.ContentId + "|*.AppImage"
}
