package main

import (
	"fmt"
	"shorturl/internal/config"
)

func main() {
	envConfig, err := config.LoadConfig(".")
	if err != nil {
		fmt.Printf("error loading config: %v\n", err)
		return
	}
	fmt.Println(envConfig)
}
