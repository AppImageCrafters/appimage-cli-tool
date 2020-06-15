package utils

// #include <stdio.h>
// #include <stdlib.h>
import "C"

import (
	"bytes"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/rainycape/dl"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
)

func MakeApplicationsDirPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	applicationsPath := filepath.Join(usr.HomeDir, "Applications")
	err = os.MkdirAll(applicationsPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	return applicationsPath, nil
}

func QueryUrl(url string) (bytes.Buffer, error) {
	resp, err := http.Get(url)
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer resp.Body.Close()

	var data bytes.Buffer
	_, err = io.Copy(&data, resp.Body)
	if err != nil {
		return bytes.Buffer{}, err
	}

	return data, nil
}

func DownloadAppImage(url string, filePath string) error {
	output, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer output.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading",
	)

	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan

		_ = resp.Body.Close()
		_ = output.Close()
		_ = os.Remove(filePath)

		os.Exit(0)
	}()

	_, err = io.Copy(io.MultiWriter(output, bar), resp.Body)
	return err
}

func UrlToTarget(target string) (string, error) {
	if strings.HasPrefix(target, "https://github.com/") {
		target, err := resolveGithubProjectTarget(target)
		return target, err
	}

	return target, nil
}

func resolveGithubProjectTarget(target string) (string, error) {
	target = target[19:]
	target_parts := strings.Split(target, "/")

	if len(target_parts) < 2 {
		return "", fmt.Errorf("missing github owner or project")
	}

	return "github:" + target_parts[0] + "/" + target_parts[1], nil
}

type DownloadLink struct {
	Name string
	Url  string
}

func PromptBinarySelection(downloadLinks []DownloadLink) (result *DownloadLink, err error) {
	if len(downloadLinks) == 1 {
		return &downloadLinks[0], nil
	}

	prompt := promptui.Select{
		Label: "Select binary to download",
		Items: downloadLinks,
		Templates: &promptui.SelectTemplates{
			Label:    "   {{ .Name }}",
			Active:   "\U00002705 {{ .Name }}",
			Inactive: "   {{ .Name }}",
			Selected: "\U00002705 {{ .Name }}"},
	}

	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return &downloadLinks[i], nil
}

func InstallAppImage(filePath string) error {
	lib, err := dl.Open("libappimage.so", 0)
	if err != nil {
		return fmt.Errorf("desktop integration not available")
	}
	defer lib.Close()

	var appimage_shall_not_be_integrated func(path *C.char) int
	err = lib.Sym("appimage_shall_not_be_integrated", &appimage_shall_not_be_integrated)
	if err != nil {
		return err
	}

	var appimage_register_in_system func(path *C.char, verbose int) int
	err = lib.Sym("appimage_register_in_system", &appimage_register_in_system)
	if err != nil {
		return err
	}

	if appimage_shall_not_be_integrated(C.CString(filePath)) != 0 {
		return nil
	}

	if appimage_register_in_system(C.CString(filePath), 1) != 0 {
		return fmt.Errorf("registration failed")
	}

	return nil
}

func MakeTargetFilePath(link *DownloadLink) (string, error) {
	applicationsPath, err := MakeApplicationsDirPath()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(applicationsPath, link.Name)
	return filePath, nil
}
