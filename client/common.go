package client

import (
	"fmt"
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
		gc, err := NewGoogleClient()
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

func HandleRequests(x XeroxApi) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		r.ParseMultipartForm(32 << 20)
		directory := x.CleanPath(strings.Join(r.PostForm[DestDir], ""))
		operation := r.PostForm[Operation]

		fmt.Println(fmt.Sprintf("Endpoint Hit: %s", operation))

		switch strings.Join(operation, "") {
		case ListDirectory:
			ListDirectoryAction(x, directory, w)
		case MakeDir:
			MakeDirectoryAction(x, directory, w)
		case PutFile:
			message, err := x.PutFile(r, directory)
			if err != nil {
				w.Write([]byte(message))
			}
		case DeleteFile:
			DeleteFileAction(x, directory, w, r)
		case RemoveDir:
			RemoveDirAction(x, directory, w)
		}
	})
}

// ListDirectory handle the list directory command
func ListDirectoryAction(x XeroxApi, directory string, w http.ResponseWriter) {
	items, err := x.ListDirectory(directory)
	if err != nil {
		//log.Println(err.Error())
		w.Write([]byte(XRXERROR))
	} else {
		w.Write([]byte(items))
	}
}

// MakeDirectory handle the make directory command
func MakeDirectoryAction(x XeroxApi, directory string, w http.ResponseWriter) {
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

// DeleteFile handle the delete file from FS
func DeleteFileAction(x XeroxApi, directory string, w http.ResponseWriter, r *http.Request) {
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

// RemoveDir handle the delete folder from FS
func RemoveDirAction(x XeroxApi, directory string, w http.ResponseWriter) {
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
