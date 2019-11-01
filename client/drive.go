package client

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Credentials struct {
	Contents Installed `json:"installed"`
}

type Installed struct {
	ClientID                string   `json:"client_id"`
	ProjectID               string   `json:"project_id"`
	AuthURI                 string   `json:"auth_uri"`
	TokenURI                string   `json:"token_uri"`
	AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
	ClientSecret            string   `json:"client_secret"`
	RedirectUris            []string `json:"redirect_uris"`
}

type Token struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.

	var tok oauth2.Token
	if filename, found := os.LookupEnv("TokenFile"); found {
		tokFile := filename
		tok, err := tokenFromFile(tokFile)
		if err != nil {
			tok = getTokenFromWeb(config)
			saveToken(tokFile, tok)
		}
	} else if accessToken, found := os.LookupEnv("AccessToken"); found {

		tokenType, err := getEnvVar("TokenType")
		if err != nil {
			return nil, err
		}

		refreshToken, err := getEnvVar("RefreshToken")
		if err != nil {
			return nil, err
		}

		expiry, err := getEnvVar("expiry")
		if err != nil {
			return nil, err
		}

		expireTime, err := time.Parse("2006-01-02T15:04:05.999999999Z", expiry)
		if err != nil {
			return nil, err
		}

		tok = oauth2.Token{AccessToken: accessToken, Expiry: expireTime, RefreshToken: refreshToken, TokenType: tokenType}
	} else {
		return nil, fmt.Errorf("google selected but no token provided")
	}

	return config.Client(context.Background(), &tok), nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	authCode = ""
	//if _, err := fmt.Scan(&authCode); err != nil {
	//	log.Fatalf("Unable to read authorization code %v", err)
	//}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getService() (*drive.Service, error) {

	var credentials []byte
	if filename, found := os.LookupEnv("credentialsFile"); found {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		credentials = b
	} else if clientId, found := os.LookupEnv("ClientId"); found {
		projectId, err := getEnvVar("ProjectID")
		if err != nil {
			return nil, err
		}

		clientSecret, err := getEnvVar("ClientSecret")
		if err != nil {
			return nil, err
		}

		uris := []string{"urn:ietf:wg:oauth:2.0:oob", "http://localhost"}
		creds := Credentials{
			Installed{
				ClientID:                clientId,
				ProjectID:               projectId,
				AuthURI:                 "https://accounts.google.com/o/oauth2/auth",
				TokenURI:                "https://oauth2.googleapis.com/token",
				AuthProviderX509CertURL: "https://www.googleapis.com/oauth2/v1/certs",
				ClientSecret:            clientSecret,
				RedirectUris:            uris,
			},
		}
		credentials, err = json.Marshal(creds)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("google selected but no credentials provided")
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(credentials, drive.DriveFileScope)

	if err != nil {
		return nil, err
	}

	client, err := getClient(config)
	if err != nil {
		fmt.Printf("Cannot create the Google Drive client: %v\n", err)
		return nil, err
	}

	service, err := drive.New(client)

	if err != nil {
		fmt.Printf("Cannot create the Google Drive service: %v\n", err)
		return nil, err
	}

	return service, err
}

func findDir(service *drive.Service, name string, parentId string) (string, error) {
	nameQuery := fmt.Sprintf("name = '%s'", name)
	parentQuery := fmt.Sprintf("'%s' in parents", parentId)
	files, err := service.Files.List().Q(fmt.Sprintf("%s and %s", nameQuery, parentQuery)).Do()

	if err != nil {
		log.Println("Could execute search: " + err.Error())
		return "", err
	}

	for _, file := range files.Files {
		fmt.Println(file.Name)
	}

	if len(files.Files) != 1 {
		return "", fmt.Errorf("%d results", len(files.Files))
	}
	return files.Files[0].Id, nil
}

func createDir(service *drive.Service, name string, parentId string) (*drive.File, error) {
	d := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentId},
	}

	file, err := service.Files.Create(d).Do()

	if err != nil {
		log.Println("Could not create dir: " + err.Error())
		return nil, err
	}

	return file, nil
}

func createFile(service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentId},
	}
	file, err := service.Files.Create(f).Media(content).Do()
	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	return file, nil
}
