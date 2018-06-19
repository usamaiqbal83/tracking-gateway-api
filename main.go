package main

import (
	"net/http"
	"strconv"
)

const PORT = 9000

func main() {
	http.HandleFunc("/", trackOpenGateway)
	err := http.ListenAndServe(":"+strconv.Itoa(PORT), nil)
	if err != nil {
		panic(err.Error())
	}
}

func trackOpenGateway(w http.ResponseWriter, r *http.Request) {
	test := "Hello World"
	w.Write([]byte(test))
}
