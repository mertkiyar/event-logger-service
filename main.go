package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Event Logger Service working good!")
	})

	fmt.Println("Service working on 8080 port")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Service not started!", err)
	}
}
