package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const ListDirectory = "ListDir"
const MakeDir = "MakeDir"
const PutFile = "PutFile"
const DeleteFile = "DeleteFile"
const RemoveDir = "RemoveDir"

const DestDir = "destDir"
const DestName = "destName"
const Operation = "theOperation"
const Sendfile = "sendfile"

const XRXNOTFOUND = "XRXNOTFOUND"
const XRXERROR = "XRXERROR"
const XRXDIREXISTS = "XRXDIREXISTS"
const XRXBADNAME = "XRXBADNAME"

type XeroxApi interface {
	ListDirectory(directory string) (string, error)
	CleanPath(directory string) string
	DeleteDir(directory string) error
	PutFile(r *http.Request, directory string) ([]byte, error)
	MakeDirectory(directory string) error
}

type Xerox struct {
	PGID int
	GID  int
}

func (x *Xerox) PutFile(r *http.Request, directory string) ([]byte, error) {
	file, _, err := r.FormFile(Sendfile)
	if err != nil {
		return []byte(XRXERROR), err
	}
	defer file.Close()

	filename := strings.Join(r.PostForm[DestName], "")
	// Create a temporary file within our temp-images directory that follows
	tempFile, err := ioutil.TempFile(directory, "")
	if err != nil {
		fmt.Println(err)
		return []byte(XRXERROR), err
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return []byte(XRXERROR), err
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	os.Rename(tempFile.Name(), fmt.Sprintf("%s/%s", directory, filename))
	if err := os.Chmod(fmt.Sprintf("%s/%s", directory, filename), 0700); err != nil {
		fmt.Println(err)
		return []byte(XRXERROR), err
	}

	if err := os.Chown(fmt.Sprintf("%s/%s", directory, filename), x.PGID, x.GID); err != nil {
		fmt.Println(err)
		return []byte(XRXERROR), err
	}

	return nil, nil
}

func (x *Xerox) ListDirectory(directory string) (string, error) {
	file, err := os.Open(directory)
	if err != nil {
		return "", err
	}
	defer file.Close()

	list, _ := file.Readdirnames(0) // 0 to read all files and folders
	directoryItems := ""
	for _, name := range list {
		fmt.Println(name)
		directoryItems = fmt.Sprintf("%s\n%s", directoryItems, name)
	}
	return directoryItems, nil
}

func (x *Xerox) MakeDirectory(directory string) error {
	err := os.Mkdir(directory, 0700)
	if err != nil {
		return err
	}
	if err := os.Chown(directory, x.PGID, x.GID); err != nil {
		return err
	}
	return nil
}

func (x *Xerox) DeleteDir(directory string) error {
	err := os.Remove(directory)
	if err != nil {
		return err
	}

	return nil
}

func (x *Xerox) CleanPath(directory string) string {
	if strings.Index(directory, "\\") >= 0 {
		directory = strings.Replace(directory, "\\", "/", -1)
	}

	for strings.Index(directory, "//") >= 0 {
		directory = strings.Replace(directory, "//", "/", -1)
	}

	return directory
}
