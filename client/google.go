package client

import (
	"bytes"
	"fmt"
	"google.golang.org/api/drive/v3"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type GoogleClient struct {
	XeroxApi
	service *drive.Service
}

// NewGoogleXeroxClient creates a new Google Drive client
func NewGoogleXeroxClient() (XeroxApi, error) {
	return NewGoogleClient()
}

// NewGoogleXeroxClient creates a new Google Drive client
func NewGoogleClient() (*GoogleClient, error) {
	service, err := getService()
	if err != nil {
		return nil, err
	}

	googleClient := GoogleClient{service: service}
	return &googleClient, nil
}

// PutFile is the function to upload a file
func (google *GoogleClient) PutFile(r *http.Request, directory string) ([]byte, error) {
	file, _, err := r.FormFile(Sendfile)
	if err != nil {
		log.Println(err)
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
		return []byte(XRXERROR), err
	}
	fmt.Printf("File '%s' successfully uploaded in '%s' directory", driveFile.Name, directory)
	if strings.Index(filename, ".pdf") > 0 && *mqttEnabled {
		err = google.submitToPubSub(&driveFile.Id, &parentId, &filename)
		if err != nil {
			log.Println("Failed to submit to mqtt")
		}
	}
	return nil, nil
}

// ListDirectory is the function to list directory
func (google *GoogleClient) ListDirectory(directory string) (string, error) {
	files, err := google.ListGoogleDirectory(directory)
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

func (google *GoogleClient) ListGoogleDirectory(directory string) (*drive.FileList, error) {
	folderId, err := google.FindDir(directory)
	if err != nil {
		return nil, err
	}
	parentQuery := fmt.Sprintf("'%s' in parents", folderId)
	files, err := google.service.Files.List().Q(parentQuery).Do()
	if err != nil {
		return nil, err
	}
	return files, nil
}

// MakeDirectory is the function to mkdir in google drive
func (google *GoogleClient) MakeDirectory(directory string) error {
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
func (google *GoogleClient) FindDir(directory string) (string, error) {
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
func (google *GoogleClient) DeleteDir(directory string) error {
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
func (google *GoogleClient) CleanPath(directory string) string {
	if strings.Index(directory, "\\") >= 0 {
		directory = strings.Replace(directory, "\\", "/", -1)
	}

	for strings.Index(directory, "//") >= 0 {
		directory = strings.Replace(directory, "//", "/", -1)
	}

	return directory
}

func (google *GoogleClient) OCRFile(fileId string, parentId string, name string) (*drive.File, error) {
	cmd := exec.Command("/usr/local/bin/ocrmypdf", "--force-ocr", "-", "-")
	gFile := google.service.Files.Get(fileId)
	dFile, err := gFile.Download()
	b, err := io.ReadAll(dFile.Body)
	if err != nil {
		log.Fatalln(err)
	}

	r := bytes.NewReader(b)
	cmd.Stdin = r
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("error: %s with output: %s", err.Error(), string(out))
		// Lets not hold up successfully storing the non OCR file
		return nil, nil
	}
	ocrReader := bytes.NewReader(out)
	focr := &drive.File{
		MimeType: "application/pdf",
		Name:     fmt.Sprintf("ocr_%s", name),
		Parents:  []string{parentId},
	}
	file, err := google.service.Files.Create(focr).Media(ocrReader).Do()
	if err != nil {
		log.Printf("error: %s with output: %s \n", err.Error(), string(out))
		// Lets not hold up successfully storing the non OCR file
		return nil, nil
	}
	log.Printf("File %s successfully uploaded in %s directory \n", focr.Name, parentId)
	err = google.service.Files.Delete(fileId).Do()
	if err != nil {
		log.Printf("Unable to delete original file: %s \n", name)
	}
	return file, nil
}
