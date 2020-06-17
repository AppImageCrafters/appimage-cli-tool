package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type RegistryEntry struct {
	Id   string
	SHA1 string
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

func (registry *Registry) Add(filePath string, id string) error {
	sha1Checksum, err := getFileSha1Checksum(filePath)
	if err != nil {
		return err
	}

	if registry.Entries == nil {
		registry.Entries = map[string]RegistryEntry{}
	}

	registry.Entries[filepath.Base(filePath)] = RegistryEntry{Id: id, SHA1: sha1Checksum}
	return nil
}

func (registry *Registry) Get(fileName string) (entry RegistryEntry, ok bool) {
	entry, ok = registry.Entries[fileName]
	return entry, ok
}

func (registry *Registry) Remove(fileName string) {
	delete(registry.Entries, fileName)
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
			_, ok := registry.Entries[f.Name()]
			if !ok {
				_ = registry.Add(filepath.Join(applicationsDir, f.Name()), "")
			}
		}
	}

	for fileName, _ := range registry.Entries {
		filePath := filepath.Join(applicationsDir, fileName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			registry.Remove(fileName)
		}
	}
}

func (registry *Registry) Lookup(id string) (string, bool) {
	for fileName, entry := range registry.Entries {
		if entry.Id == id {
			return fileName, true
		}
	}

	return "", false
}

func getFileSha1Checksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	sha1Checksum := sha1.New()
	_, err = io.Copy(sha1Checksum, file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sha1Checksum.Sum(nil)), nil
}

func makeRegistryFilePath() (string, error) {
	applicationsPath, err := MakeApplicationsDirPath()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(applicationsPath, ".registry.json")
	return filePath, nil
}
