package commands

import (
	"fmt"
	"github.com/tidwall/gjson"
)

type SearchCmd struct {
	Query string `arg name:"query" help:"Query be used in the search." type:"string"`
}

func (r *SearchCmd) Run(*Context) error {
	jsonData, err := queryUrl("https://www.pling.com/json/search/p/" + r.Query + "/s/AppImageHub.com")
	if err != nil {
		return err
	}

	jsonParser := gjson.Parse(jsonData.String())

	products := jsonParser.Get(`#(title="Products").values`)
	for _, product := range products.Array() {
		id := product.Get(`project_id`).String()
		title := product.Get(`title`).String()
		category := product.Get(`cat_title`).String()
		username := product.Get(`username`).String()
		fmt.Printf("appimagehub:%s %s (%s) by %s\n", id, title, category, username)
	}

	return nil
}
