package utils

import (
	"bytes"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
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
