package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// exit println 'msgs' and exit with 'code'
func exit(code int, msgs ...interface{}) {
	fmt.Println(msgs...)
	os.Exit(code)
}

// exitf prints formatted 'msgs' using 'tpl' and exit with 'code'
func exitf(code int, tpl string, msgs ...interface{}) {
	fmt.Printf(tpl, msgs...)
	os.Exit(code)
}

// showDebug prints values if debug enabled
func showDebug(v ...interface{}) {
	if *debug {
		log.Println(v...)
	}
}

// showDebugF prints foormated values if debug enabled
func showDebugF(format string, v ...interface{}) {
	if !*debug {
		return
	}

	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}

	log.Printf(format, v...)
}

// isROOT checks if the current user is ROOT
// works only on unix/linux systems.
func isROOT() bool {
	if os.Geteuid() == 0 {
		return true
	}
	return false
}
