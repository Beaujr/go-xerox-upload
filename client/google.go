package client

import (
	"fmt"
	"google.golang.org/api/drive/v3"
	"net/http"
	"os"
	"strings"
	"time"
)

type googleClient struct {
	XeroxApi
	service *drive.Service
}

func NewGoogleClient() (XeroxApi, error) {
	service, err := getService()
	if err != nil {
		return nil, err
	}

	googleClient := googleClient{service: service}
	return &googleClient, nil
}

func (google *googleClient) PutFile(r *http.Request, directory string) ([]byte, error) {
	file, _, err := r.FormFile(Sendfile)
	if err != nil {
		return []byte(XRXERROR), err
	}
	defer file.Close()
	filename := strings.Join(r.PostForm[DestName], "")
	filename = fmt.Sprintf("%s_%s", time.Now().Format("20060102150405"), filename)

	parentId, err := google.FindDir(directory)
	if err != nil {
		return []byte(XRXERROR), err
	}

	driveFile, err := createFile(google.service, filename, "application/pdf", file, parentId)
	if err != nil {
		panic(fmt.Sprintf("Could not create file: %v\n", err))
	}

	fmt.Printf("File '%s' successfully uploaded in '%s' directory", driveFile.Name, directory)

	return nil, nil
}

func (google *googleClient) ListDirectory(directory string) (string, error) {
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

func (google *googleClient) MakeDirectory(directory string) error {
	directories := strings.Split(directory, "/")
	root := "root"
	for _, dir := range directories {
		if dir == "" {
			continue
		}
		folderId, err := findDir(google.service, dir, root)
		if err != nil {
			fmt.Println(err.Error())
			if strings.Compare(err.Error(), "0 results") == 0 {
				folder, err := createDir(google.service, dir, root)
				if err != nil {
					return err
				}
				folderId = folder.Id
			} else {
				return err
			}
		}
		root = folderId
	}
	return nil
}

func (google *googleClient) FindDir(directory string) (string, error) {
	directories := strings.Split(directory, "/")
	root := "root"
	for _, dir := range directories {
		if dir == "" {
			continue
		}
		folderId, err := findDir(google.service, dir, root)
		if err != nil {
			fmt.Println(err.Error())
			if strings.Compare(err.Error(), "0 results") == 0 {
				folder, err := createDir(google.service, dir, root)
				if err != nil {
					return "", err
				}
				folderId = folder.Id
			} else {
				return "", err
			}
		}
		root = folderId
	}

	return root, nil
}

func (google *googleClient) DeleteDir(directory string) error {
	err := os.Remove(directory)
	if err != nil {
		return err
	}

	return nil
}

func (google *googleClient) CleanPath(directory string) string {
	if strings.Index(directory, "\\") >= 0 {
		directory = strings.Replace(directory, "\\", "/", -1)
	}

	for strings.Index(directory, "//") >= 0 {
		directory = strings.Replace(directory, "//", "/", -1)
	}

	return directory
}
