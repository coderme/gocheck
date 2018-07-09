package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// Result defines the result of
// checking certain URL
type Result struct {
	Code    string        `json:"code"`
	URL     string        `json:"url"`
	Time    time.Time     `json:"time"`
	Elapsed time.Duration `json:"took"`

	ErrorConnect  bool `json:"error_connect"`
	ErrorClient   bool `json:"error_client"`
	ErrorServer   bool `json:"error_server"`
	ErrorRedirect bool `json:"error_redirect"`
}

// PrintText prints result in an easy format
func (r *Result) PrintText() {
	fmt.Printf("%-9s %s\n", r.Code, r.URL)
}

// PrintJSON prints json value of the result
func (r *Result) PrintJSON() {
	j, err := json.Marshal(r)
	if err != nil {
		exit(1, "Failed json.Marshal", err)
	}
	fmt.Println(string(j))
}
