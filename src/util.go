package main

import (
	"fmt"
	"os"
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

// isROOT checks if the current user is ROOT
// works only on unix/linux systems.
func isROOT() bool {
	if os.Geteuid() == 0 {
		return true
	}
	return false
}
