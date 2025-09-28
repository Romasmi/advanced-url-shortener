package main

import (
	"fmt"
	"shorturl/internal/application"
)

func main() {
	app := &application.App{}
	err := app.InitApp("../../")
	if err != nil {
		fmt.Printf("error while app initialization: %v", err)
		return
	}
	defer app.OnStop()

}
