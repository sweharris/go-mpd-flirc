package main

import (
	"fmt"
	"os"
)

func die(v ...any) {
	fmt.Fprintln(os.Stderr, "\nError:")
	fmt.Fprintf(os.Stderr, "  %v\n", v...)
	os.Exit(1)
}


func ignore_it() bool {
        n := current_window()
        if n == "Kodi" || n == "" {
                return true
        }
        return false
}

