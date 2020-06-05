package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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

func (registry *Registry) Set(fileName string, id string) error {
	sha1Checksum, err := getFileSha1Checksum(fileName)
	if err != nil {
		return err
	}

	if registry.Entries == nil {
		registry.Entries = map[string]RegistryEntry{}
	}

	registry.Entries[fileName] = RegistryEntry{Id: id, SHA1: sha1Checksum}
	return nil
}

func (registry *Registry) Get(fileName string) RegistryEntry {
	return registry.Entries[fileName]
}

func (registry *Registry) Unset(fileName string) {

	delete(registry.Entries, fileName)
}

func getFileSha1Checksum(fileName string) (string, error) {
	applicationsDir, err := MakeApplicationsDirPath()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(applicationsDir, fileName)

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
