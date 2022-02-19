package client

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// ListDirectory is the Payload value from the Printer to List Directory Values to avoid filename collisions
const ListDirectory = "ListDir"

// MakeDir is the Payload value from the Printer to Make a Directory if it doesn't exit
const MakeDir = "MakeDir"

// PutFile is the Payload value from the Printer to upload a file
const PutFile = "PutFile"

// DeleteFile is the Payload value from the Printer to Delete a file normally a *LCK file
const DeleteFile = "DeleteFile"

// RemoveDir is the Payload value from the Printer to Delete a directory
const RemoveDir = "RemoveDir"

// DestDir is the Payload value from the Printer for the files dir on the filesystem
const DestDir = "destDir"

// DestName is the Payload value from the Printer for the filename  on the filesystem
const DestName = "destName"

// Operation is the Payload field from the Printer for the the operation to happen, ListDir, MakeDir etc
const Operation = "theOperation"

// Sendfile is the Payload field from the Printer for the the file itself
const Sendfile = "sendfile"

// XRXNOTFOUND is the not found error message
const XRXNOTFOUND = "XRXNOTFOUND"

// XRXERROR is the default error message
const XRXERROR = "XRXERROR"

// XRXDIREXISTS is the directory exists already
const XRXDIREXISTS = "XRXDIREXISTS"

// XRXBADNAME is the filename is bad due to to FS constraints
const XRXBADNAME = "XRXBADNAME"

// XeroxApi Interface for all Printer to Server Interactions
type XeroxApi interface {
	ListDirectory(directory string) (string, error)
	CleanPath(directory string) string
	DeleteDir(directory string) error
	PutFile(r *http.Request, directory string) ([]byte, error)
	MakeDirectory(directory string) error
}

// NewClient generates a new generic client for uploading
func NewClient() (XeroxApi, error) {
	_, found := os.LookupEnv("google")
	var x XeroxApi
	if found {
		gc, err := NewGoogleXeroxClient()
		if err != nil {
			return nil, err
		}
		x = gc
	} else {
		pgId, err := getEnvVar("PGID")
		if err != nil {
			return nil, err
		}

		userId, err := strconv.Atoi(pgId)
		if err != nil {
			return nil, err
		}

		gID, err := getEnvVar("GID")
		if err != nil {
			return nil, err
		}

		groupIp, err := strconv.Atoi(gID)
		if err != nil {
			return nil, err
		}
		x = NewFileSystemClient(userId, groupIp)
	}
	return x, nil
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
func HandleOCRRequest() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, found := os.LookupEnv("cloudrun")
		if found {
			user, pass, ok := r.BasicAuth()
			if !ok || user != "xerox" && pass != "5.[H]/_qpgq39[{t" {
				w.WriteHeader(401)
				w.Write([]byte("Unauthorized.\n"))
				return
			}
		}
		x, err := NewGoogleClient()
		if err != err {
			w.WriteHeader(401)
			w.Write([]byte("Unauthorized.\n"))
			return
		}

		files, err := x.ListGoogleDirectory("/mail/shared/2022")
		if err != err {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		parentId, err := x.FindDir("/mail/shared/2022")
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		for _, file := range files.Files {
			if file.MimeType == "application/pdf" && strings.Index(file.Name, "ocr") < 0 {
				log.Printf("OCRing %s\n", file.Name)
				_, err := x.OCRFile(file.Id, parentId, file.Name)
				if err != nil {
					log.Printf("Error occured for file %s, %s\n", file.Name, err.Error())
				}
			}
		}
		w.WriteHeader(200)
		return
	})
}

// HandleRequests takes the XeroxApi and handles all the List, Del, Remove, Put actions
func HandleRequests(x XeroxApi) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, found := os.LookupEnv("cloudrun")
		if found {
			user, pass, ok := r.BasicAuth()
			if !ok || user != "xerox" && pass != "5.[H]/_qpgq39[{t" {
				w.WriteHeader(401)
				w.Write([]byte("Unauthorized.\n"))
				return
			}
		}

		w.WriteHeader(200)
		r.ParseMultipartForm(32 << 20)
		directory := x.CleanPath(strings.Join(r.PostForm[DestDir], ""))
		operation := r.PostForm[Operation]
		fmt.Println(fmt.Sprintf("Endpoint Hit: %s", operation))

		switch strings.Join(operation, "") {
		case ListDirectory:
			listDirectoryAction(x, directory, w)
		case MakeDir:
			makeDirectoryAction(x, directory, w)
		case PutFile:
			message, err := x.PutFile(r, directory)
			if err != nil {
				w.Write([]byte(message))
			}
		case DeleteFile:
			deleteFileAction(x, directory, w, r)
		case RemoveDir:
			removeDirAction(x, directory, w)
		default:
			w.WriteHeader(200)
		}
	})
}

// listDirectoryAction handle the list directory command
func listDirectoryAction(x XeroxApi, directory string, w http.ResponseWriter) {
	items, err := x.ListDirectory(directory)
	if err != nil {
		log.Println(err.Error())
		w.Write([]byte(XRXERROR))
	} else {
		w.Write([]byte(items))
	}
}

// makeDirectoryAction handle the make directory command
func makeDirectoryAction(x XeroxApi, directory string, w http.ResponseWriter) {
	err := x.MakeDirectory(directory)
	if err != nil {
		//   XRXBADNAME if the name is empty.
		//   XRXDIREXISTS if the directory already exists.
		//   XRXERROR if the directory cannot be created.
		switch err.(*os.PathError).Err {
		case syscall.EEXIST:
			w.Write([]byte(XRXDIREXISTS))
		case syscall.ENOENT:
			w.Write([]byte(XRXBADNAME))
		default:
			w.Write([]byte(XRXERROR))
		}
	}
}

// deleteFileAction handle the delete file from FS
func deleteFileAction(x XeroxApi, directory string, w http.ResponseWriter, r *http.Request) {
	//   XRXNOTFOUND if the requested file isn't found.
	//   XRXERROR the file cannot be deleted.
	destinationName := r.PostForm[DestName]
	filename := strings.Join(destinationName, "")
	if strings.Join(destinationName, "") != "" {
		directory = fmt.Sprintf("%s%s", directory, filename)
	}
	err := x.DeleteDir(directory)
	if err != nil {
		switch err.(*os.PathError).Err {
		case syscall.ENOENT:
			w.Write([]byte(XRXNOTFOUND))
		default:
			w.Write([]byte(XRXERROR))
		}
	}
}

// removeDirAction handle the delete folder from FS
func removeDirAction(x XeroxApi, directory string, w http.ResponseWriter) {
	//   XRXBADNAME if the requested file isn't of the correct type or the name is empty.
	//   XRXNOTFOUND if the requested file isn't found.
	//   XRXERROR the file cannot be deleted.
	err := x.DeleteDir(directory)
	if err != nil {
		switch err.(*os.PathError).Err {
		case syscall.ENOENT:
			w.Write([]byte(XRXNOTFOUND))
		default:
			w.Write([]byte(XRXERROR))
		}
	}
}
