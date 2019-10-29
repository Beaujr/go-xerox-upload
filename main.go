package main

import (
	"fmt"
	xclient "github.com/beaujr/go-xerox-upload/client"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/upload", upload)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Xerox - Go server")
	handleRequests()
}

func upload(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	r.ParseMultipartForm(32 << 20)

	var x xclient.XeroxApi
	x, err := xclient.NewClient()
	if err != nil {
		log.Println(err.Error())
		w.Write([]byte(xclient.XRXERROR))
		return
	}

	directory := x.CleanPath(strings.Join(r.PostForm[xclient.DestDir], ""))
	operation := r.PostForm[xclient.Operation]

	fmt.Println(fmt.Sprintf("Endpoint Hit: %s", operation))

	switch strings.Join(operation, "") {
	case xclient.ListDirectory:
		ListDirectory(x, directory, w, r)
	case xclient.MakeDir:
		MakeDirectory(x, directory, w, r)
	case xclient.PutFile:
		message, err := x.PutFile(r, directory)
		if err != nil {
			w.Write([]byte(message))
		}
	case xclient.DeleteFile:
		DeleteFile(x, directory, w, r)
	case xclient.RemoveDir:
		RemoveDir(x, directory, w, r)
	}
}

// ListDirectory handle the list directory command
func ListDirectory(x xclient.XeroxApi, directory string, w http.ResponseWriter, r *http.Request) {
	items, err := x.ListDirectory(directory)
	if err != nil {
		log.Println(err.Error())
		w.Write([]byte(xclient.XRXERROR))

	} else {
		w.Write([]byte(items))
	}
}

// MakeDirectory handle the make directory command
func MakeDirectory(x xclient.XeroxApi, directory string, w http.ResponseWriter, r *http.Request) {
	err := x.MakeDirectory(directory)
	if err != nil {
		//   XRXBADNAME if the name is empty.
		//   XRXDIREXISTS if the directory already exists.
		//   XRXERROR if the directory cannot be created.
		switch err.(*os.PathError).Err {
		case syscall.EEXIST:
			w.Write([]byte(xclient.XRXDIREXISTS))
		case syscall.ENOENT:
			w.Write([]byte(xclient.XRXBADNAME))
		default:
			w.Write([]byte(xclient.XRXERROR))
		}
	}
}

// DeleteFile handle the delete file from FS
func DeleteFile(x xclient.XeroxApi, directory string, w http.ResponseWriter, r *http.Request) {
	//   XRXNOTFOUND if the requested file isn't found.
	//   XRXERROR the file cannot be deleted.
	destinationName := r.PostForm[xclient.DestName]
	filename := strings.Join(destinationName, "")
	if strings.Join(destinationName, "") != "" {
		directory = fmt.Sprintf("%s%s", directory, filename)
	}
	err := x.DeleteDir(directory)
	if err != nil {
		switch err.(*os.PathError).Err {
		case syscall.ENOENT:
			w.Write([]byte(xclient.XRXNOTFOUND))
		default:
			w.Write([]byte(xclient.XRXERROR))
		}
	}
}

// RemoveDir handle the delete folder from FS
func RemoveDir(x xclient.XeroxApi, directory string, w http.ResponseWriter, r *http.Request) {
	//   XRXBADNAME if the requested file isn't of the correct type or the name is empty.
	//   XRXNOTFOUND if the requested file isn't found.
	//   XRXERROR the file cannot be deleted.
	err := x.DeleteDir(directory)
	if err != nil {
		switch err.(*os.PathError).Err {
		case syscall.ENOENT:
			w.Write([]byte(xclient.XRXNOTFOUND))
		default:
			w.Write([]byte(xclient.XRXERROR))
		}
	}
}
