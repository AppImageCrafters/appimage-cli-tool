package utils

import (
	"encoding/json"
	updateUtils "github.com/AppImageCrafters/appimage-update/util"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type RegistryEntry struct {
	Repo       string
	FileSha1   string
	AppName    string
	AppVersion string
	FilePath   string
	UpdateInfo string
}

type Registry struct {
	Entries map[string]RegistryEntry
}

func OpenRegistry() (registry *Registry, err error) {
	path, err := makeRegistryFilePath()
	if err != nil {
		return
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return &Registry{Entries: map[string]RegistryEntry{}}, nil
	}

	err = json.Unmarshal(data, &registry)
	if err != nil {
		return
	}

	return
}

func (registry *Registry) Close() error {
	path, err := makeRegistryFilePath()
	if err != nil {
		return err
	}

	blob, err := json.Marshal(registry)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, blob, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (registry *Registry) Add(entry RegistryEntry) error {
	registry.Entries[entry.FilePath] = entry
	return nil
}

func (registry *Registry) Remove(filePath string) {
	delete(registry.Entries, filePath)
}

func (registry *Registry) Update() {
	applicationsDir, err := MakeApplicationsDirPath()
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(applicationsDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".AppImage") {
			filePath := filepath.Join(applicationsDir, f.Name())
			_, ok := registry.Entries[filePath]
			if !ok {
				entry := registry.createEntryFromFile(filePath)
				_ = registry.Add(entry)
			}
		}
	}

	for filePath, _ := range registry.Entries {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			registry.Remove(filePath)
		}
	}
}

func (registry *Registry) addFile(filePath string) {
	entry := registry.createEntryFromFile(filePath)
	_ = registry.Add(entry)
}

func (registry *Registry) createEntryFromFile(filePath string) RegistryEntry {
	fileSha1, _ := GetFileSHA1(filePath)
	updateInfo, _ := updateUtils.ReadUpdateInfo(filePath)
	entry := RegistryEntry{
		Repo:       "",
		FileSha1:   fileSha1,
		AppName:    "",
		AppVersion: "",
		FilePath:   filePath,
		UpdateInfo: updateInfo,
	}
	return entry
}

func (registry *Registry) Lookup(target string) (RegistryEntry, bool) {
	applicationsDir, _ := MakeApplicationsDirPath()
	possibleFullPath := filepath.Join(applicationsDir, target)

	for _, entry := range registry.Entries {
		if entry.FileSha1 == target || entry.FilePath == target ||
			entry.FilePath == possibleFullPath || entry.Repo == target {
			return entry, true
		}
	}

	if IsAppImageFile(target) {
		entry := registry.createEntryFromFile(target)
		_ = registry.Add(entry)

		return entry, true
	} else {
		if IsAppImageFile(possibleFullPath) {
			entry := registry.createEntryFromFile(target)
			_ = registry.Add(entry)

			return entry, true
		}
	}

	return RegistryEntry{}, false
}

func makeRegistryFilePath() (string, error) {
	applicationsPath, err := MakeApplicationsDirPath()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(applicationsPath, ".registry.json")
	return filePath, nil
}
