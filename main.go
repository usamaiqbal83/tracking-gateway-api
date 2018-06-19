package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const PORT = 9000
const apiURL = "https://staging.leadfuze.com"

//https://lfclk.co/trk/o?t=ODY2NXwyNTYyMnwxMDY4Mjc3MQ==

func main() {
	http.HandleFunc("/trk/o", trackOpenGateway)
	fmt.Println("Server Starting listening at : " +strconv.Itoa(PORT))
	err := http.ListenAndServe(":"+strconv.Itoa(PORT), nil)
	if err != nil {
		panic(err.Error())
	}
}

func trackOpenGateway(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		token := r.URL.Query().Get("t")
		if token != ""{
			url := fmt.Sprintf(apiURL+"/trk/o?t=%s",token)
			fmt.Println("URL : ", url)
			res, err := http.Get(url)
			if err != nil {
				fmt.Println("Error : ", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				defer res.Body.Close()
				contents, err := ioutil.ReadAll(res.Body)
				if err != nil {
					fmt.Println("Error : ", err.Error())
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "image/png")
				w.Header().Set("Cache-Control", "no-cache, max-age=0")
				w.Write(contents)
				w.WriteHeader(res.StatusCode)
				//w.WriteHeader(http.StatusOK)
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internet Server Error"))
	}
}
