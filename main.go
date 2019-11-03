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
		appengine.Main()
	} else {
		myRouter := mux.NewRouter().StrictSlash(true)
		myRouter.Handle("/upload", xclient.HandleRequests(x))
		http.ListenAndServe(":10000", myRouter)
	}
}
