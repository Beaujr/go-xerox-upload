package main

import (
	"fmt"
	xclient "github.com/beaujr/go-xerox-upload/client"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
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

	pgId, err := getEnvVar("PGID")
	if err != nil {
		log.Panic(err)
	}

	userId, err := strconv.Atoi(pgId)
	if err != nil {
		log.Panic(err)
	}

	gID, err := getEnvVar("GID")
	if err != nil {
		log.Panic(err)
	}

	groupIp, err := strconv.Atoi(gID)
	if err != nil {
		log.Panic(err)
	}

	_, found := os.LookupEnv("google")
	var x xclient.XeroxApi

	if found {
		x, err = xclient.NewGoogleClient(userId, groupIp)
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(xclient.XRXERROR))
			return
		}
	} else {
		x = &xclient.Xerox{PGID: userId, GID: groupIp}
	}

	directory := x.CleanPath(strings.Join(r.PostForm[xclient.DestDir], ""))

	operation := r.PostForm[xclient.Operation]

	fmt.Println(fmt.Sprintf("Endpoint Hit: %s", operation))

	switch strings.Join(operation, "") {
	case xclient.ListDirectory:
		items, err := x.ListDirectory(directory)
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(xclient.XRXERROR))

		} else {
			w.Write([]byte(items))
		}
	case xclient.MakeDir:
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
	case xclient.PutFile:
		message, err := x.PutFile(r, directory)
		if err != nil {
			w.Write([]byte(message))
		}

	case xclient.DeleteFile:
		//   XRXNOTFOUND if the requested file isn't found.
		//   XRXERROR the file cannot be deleted.
		destinationName := r.PostForm[xclient.DestName]
		filename := strings.Join(destinationName, "")
		if strings.Join(destinationName, "") != "" {
			directory = fmt.Sprintf("%s%s", directory, filename)
		}
		err := os.Remove(directory)
		if err != nil {
			switch err.(*os.PathError).Err {
			case syscall.ENOENT:
				w.Write([]byte(xclient.XRXNOTFOUND))
			default:
				w.Write([]byte(xclient.XRXERROR))
			}
		}
	case xclient.RemoveDir:
		//   XRXBADNAME if the requested file isn't of the correct type or the name is empty.
		//   XRXNOTFOUND if the requested file isn't found.
		//   XRXERROR the file cannot be deleted.
		err := os.Remove(directory)
		if err != nil {
			switch err.(*os.PathError).Err {
			case syscall.ENOENT:
				w.Write([]byte(xclient.XRXNOTFOUND))
			default:
				w.Write([]byte(xclient.XRXERROR))
			}
		}
	}

}

func getEnvVar(name string) (string, error) {
	v, found := os.LookupEnv(name)
	if !found {
		return "", fmt.Errorf("%s must be set", name)
	}
	if len(v) == 0 {
		return "", fmt.Errorf("%s must not be empty", name)
	}
	return v, nil
}
