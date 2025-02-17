package main

import (
	"fmt"

	"github.com/vehsamrak/go-version-parser/internal/application"
	"github.com/vehsamrak/go-version-parser/internal/markdown"
)

func main() {
	app := application.NewApplication(markdown.NewMarkdownDrawer())

	err := app.Execute(
		application.Input{
			IgnoreVendors: true,
			OutputType:    application.OutputTypeMarkdown,
			RootDirectory: ".",
		},
	)
	if err != nil {
		fmt.Printf("error: %s", err)
	}
}
