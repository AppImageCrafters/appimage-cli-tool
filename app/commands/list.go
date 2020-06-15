package commands

import (
	"appimage-manager/app/utils"
	"bytes"
	"github.com/juju/ansiterm"
	"os"
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
	_, _ = tabWriter.Write([]byte("Target\t File Name\t SHA1\n"))
	_, _ = tabWriter.Write([]byte("--\t ---------\t ----\n"))

	tabWriter.SetForeground(ansiterm.DarkGray)

	for fileName, v := range registry.Entries {
		_, _ = tabWriter.Write([]byte(v.Id + "\t " + fileName + "\t " + v.SHA1 + "\n"))
	}
	_ = tabWriter.Flush()
	_, _ = os.Stdout.Write(buf.Bytes())
	return nil
}
