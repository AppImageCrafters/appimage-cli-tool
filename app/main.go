package main

import (
	"appimage-installer/app/commands"
	"github.com/alecthomas/kong"
)

var cli struct {
	Debug bool `help:"Enable debug mode."`

	Search  commands.SearchCmd  `cmd help:"Search applications in the store."`
	Install commands.InstallCmd `cmd help:"Install application."`
	List    commands.ListCmd    `cmd help:"List installed applications."`
}

func main() {
	ctx := kong.Parse(&cli)
	// Call the Run() method of the selected parsed command.
	err := ctx.Run(&commands.Context{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)

}
