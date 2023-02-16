package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestShawg(t *testing.T) {
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":9090", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello")
}
