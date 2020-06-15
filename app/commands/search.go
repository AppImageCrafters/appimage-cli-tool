package commands

import (
	"appimage-manager/app/utils"
	"bytes"
	"fmt"
	"github.com/juju/ansiterm"
	"github.com/tidwall/gjson"
	"os"
)

type SearchCmd struct {
	Query string `arg name:"query" help:"Query be used in the search." type:"string"`
}

func (r *SearchCmd) Run(*Context) error {
	jsonData, err := utils.QueryUrl("https://www.pling.com/json/search/p/" + r.Query + "/s/AppImageHub.com")
	if err != nil {
		return err
	}

	jsonParser := gjson.Parse(jsonData.String())

	products := jsonParser.Get(`#(title="Products").values`)

	var buf bytes.Buffer
	tabWriter := ansiterm.NewTabWriter(&buf, 20, 4, 0, ' ', 0)
	tabWriter.SetColorCapable(true)

	tabWriter.SetForeground(ansiterm.Green)
	_, _ = tabWriter.Write([]byte("Target\t Name\t Category\t Publisher\n"))
	_, _ = tabWriter.Write([]byte("--\t ----\t --------\t ---------\n"))

	tabWriter.SetForeground(ansiterm.DarkGray)

	for _, product := range products.Array() {
		id := product.Get(`project_id`).String()
		title := product.Get(`title`).String()
		category := product.Get(`cat_title`).String()
		username := product.Get(`username`).String()

		fmt.Fprintf(tabWriter, "appimagehub:%s\t %s\t %s\t %s\n", id, title, category, username)
	}

	_ = tabWriter.Flush()
	_, _ = os.Stdout.Write(buf.Bytes())

	return nil
}
