package main

import (
	"fmt"

	"github.com/Romasmi/advanced-url-shortener/internal/application"
)

func main() {
	app := &application.App{}
	err := app.InitApp("../../")
	if err != nil {
		fmt.Printf("error while app initialization: %v", err)
		return
	}
	defer app.OnStop()

	app.Run()
}
