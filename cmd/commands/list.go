package commands

import (
	"appimage-cli-tool/internal/utils"
	"bytes"
	"github.com/juju/ansiterm"
	"os"
	"path/filepath"
	"sort"
)

type ListCmd struct {
}

func (r *ListCmd) Run(*Context) error {
	registry, err := utils.OpenRegistry()
	if err != nil {
		return err
	}
	defer registry.Close()

	registry.Update()
	var buf bytes.Buffer
	tabWriter := ansiterm.NewTabWriter(&buf, 20, 4, 0, ' ', 0)
	tabWriter.SetColorCapable(true)

	tabWriter.SetForeground(ansiterm.Green)
	_, _ = tabWriter.Write([]byte("Repo\t File Name\t SHA1\n"))
	_, _ = tabWriter.Write([]byte("----\t ---------\t ----\n"))

	tabWriter.SetForeground(ansiterm.DarkGray)

	var lines [][]string
	for fileName, v := range registry.Entries {
		line := []string{v.Repo, filepath.Base(fileName), v.FileSha1}
		lines = append(lines, line)
	}
	sort.Slice(lines, func(i int, j int) bool {
		return lines[i][1] < lines[j][1]
	})

	for _, line := range lines {
		_, _ = tabWriter.Write([]byte(line[0] + "\t " + line[1] + "\t " + line[2] + "\n"))
	}
	_ = tabWriter.Flush()
	_, _ = os.Stdout.Write(buf.Bytes())
	return nil
}
