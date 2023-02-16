package main

import (
	"fmt"
	"net/http"
	"shawf/shawg/shawg"
)

func main() {
	r := shawg.New()
	r.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello,my man")
	})

	r.GET("/hi", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q]=%q\n", k, v)
		}
	})
	r.Run(":3023")
}
