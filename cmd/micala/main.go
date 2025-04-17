package main

import (
	"github.com/yangszwei/koala/internal/registry"
)

func main() {
	app := registry.NewApp()

	if err := app.Init(); err != nil {
		panic(err)
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}
