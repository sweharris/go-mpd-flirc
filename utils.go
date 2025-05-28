package main

import (
	"log"
	"os"
)

func die(v ...any) {
	log.Println("\nError:")
	log.Printf("  %v\n", v...)
	os.Exit(1)
}

func ignore_it() bool {
	n := current_window()
	if n == "Kodi" || n == "" {
		return true
	}
	return false
}
