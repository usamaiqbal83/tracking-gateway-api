package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

const PORT = 9000
const apiURL = "https://staging.leadfuze.com"

//https://lfclk.co/trk/o?t=ODY2NXwyNTYyMnwxMDY4Mjc3MQ==

func main() {
	http.HandleFunc("/trk/o", trackOpenGateway)
	logrus.Println("Server Starting listening at : " +strconv.Itoa(PORT))
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
			logrus.Println("URL : ", url)
			res, err := http.Get(url)
			if err != nil {
				logrus.Println("Error : ", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				defer res.Body.Close()
				contents, err := ioutil.ReadAll(res.Body)
				if err != nil {
					logrus.Println("Error : ", err.Error())
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "image/png")
				w.Header().Set("Cache-Control", "no-cache, max-age=0")
				w.Write(contents)
				w.WriteHeader(res.StatusCode)
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

func InitLogger() {

	env := os.Getenv("environment")
	isLocalHost := env == "local"

	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	if isLocalHost {
		logrus.SetOutput(os.Stdout)
	}

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.InfoLevel)

	if !isLocalHost {
		// configure file system hook
		configureLocalFileSystemHook()
	}
}

func configureLocalFileSystemHook() {
	//path := "./log/go.log"
	path := "/var/log/tracker.log"

	rLogs, err := rotatelogs.New(
		path+".%Y_%m_%d_%H_%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(30*86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(86400)*time.Second),
	)

	if err != nil {
		logrus.Println("Local file system hook initialize fail")
		return
	}

	logrus.AddHook(lfshook.NewHook(lfshook.WriterMap{
		logrus.InfoLevel:  rLogs,
		logrus.ErrorLevel: rLogs,
	}, &logrus.JSONFormatter{}))
}
