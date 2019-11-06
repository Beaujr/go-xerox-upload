package main

import (
	"fmt"
	xclient "github.com/beaujr/go-xerox-upload/client"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Xerox - Go server")
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
		appengine.Main()
	} else {
		myRouter := mux.NewRouter().StrictSlash(true)
		myRouter.Handle("/upload", xclient.HandleRequests(x))
		http.ListenAndServe(":10000", myRouter)
	}
}
