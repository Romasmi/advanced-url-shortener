package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	fmt.Printf("error checking file %s: %v\n", filename, err)
	return false
}

func FatalIfFileDoesNotExists(filename string, error string) {
	if !FileExists(filename) {
		log.Fatalln(error)
	}
}
