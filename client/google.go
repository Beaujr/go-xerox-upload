package client

import (
	"fmt"
	"google.golang.org/api/drive/v3"
	"net/http"
	"strings"
	"time"
)

type googleClient struct {
	XeroxApi
	service *drive.Service
}

// NewGoogleClient creates a new Google Drive client
func NewGoogleClient() (XeroxApi, error) {
	service, err := getService()
	if err != nil {
		return nil, err
	}

	googleClient := googleClient{service: service}
	return &googleClient, nil
}

// PutFile is the function to upload a file
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

// ListDirectory is the function to list directory
func (google *googleClient) ListDirectory(directory string) (string, error) {
	folderId, err := google.FindDir(directory)
	if err != nil {
		return "", err
	}
	parentQuery := fmt.Sprintf("'%s' in parents", folderId)
	files, err := google.service.Files.List().Q(parentQuery).Do()
	if err != nil {
		return "", err
	}
	directoryItems := ""
	for _, name := range files.Files {
		fmt.Println(name.Name)
		directoryItems = fmt.Sprintf("%s\n%s", directoryItems, name.Name)
	}
	return directoryItems, nil
}

// MakeDirectory is the function to mkdir in google drive
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

// FindDir is the function to find dir in Google Drive
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

// DeleteDir is the function to rm -rf dir
func (google *googleClient) DeleteDir(directory string) error {
	dirId, err := google.FindDir(directory)
	if err != nil {
		if strings.Compare(err.Error(), "0 results") == 0 {
			return nil
		}
		return err
	}
	err = google.service.Files.Delete(dirId).Do()
	if err != nil {
		return err
	}

	return nil
}

// CleanPath is the function to clean the path that the printer sends
func (google *googleClient) CleanPath(directory string) string {
	if strings.Index(directory, "\\") >= 0 {
		directory = strings.Replace(directory, "\\", "/", -1)
	}

	for strings.Index(directory, "//") >= 0 {
		directory = strings.Replace(directory, "//", "/", -1)
	}

	return directory
}
