package main

import (
	"flag"
	"fmt"
	xclient "github.com/beaujr/go-xerox-upload/client"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Xerox - Go server")
	flag.Parse()
	x, err := xclient.NewClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, appeng := os.LookupEnv("appengine"); appeng {
		http.Handle("/upload", xclient.HandleRequests(x))
		http.Handle("/_ah/stop", http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}))
		http.Handle("/_ah/start", http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}))
	} else {
		myRouter := mux.NewRouter().StrictSlash(true)
		myRouter.Handle("/upload", xclient.HandleRequests(x))
		myRouter.Handle("/ocr", xclient.HandleOCRRequest())
		port := os.Getenv("PORT")
		if port == "" {
			port = "10000"
		}
		if err := http.ListenAndServe(fmt.Sprintf(":%s", port), myRouter); err != nil {
			log.Panic(err)
		}
	}
}
