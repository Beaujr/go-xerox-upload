package client

import (
	"os"
	"testing"
)

func TestGetServiceResultError(t *testing.T) {
	//rec := httptest.NewRecorder()
	//req, _ := http.NewRequest(http.MethodPost, "/v1/user/FakeUserId/entitlement", bytes.NewBuffer([]byte(validEntitlementJson)))
	//env := Env{db: &mockSuccessAdminDB{}}
	//
	//test_helpers.InvokeHandler(http.Handler(createEntitlement(&env)), entitlementEndpointPath, rec, req)
	os.Setenv("google", "true")
	_, err := getService()
	expected := "google selected but no credentials provided"
	obtained := err.Error()
	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
}

func TestGetServiceResultEnv(t *testing.T) {
	os.Setenv("ClientId", "ClientId")
	_, err := getService()
	expected := "ProjectID must be set"
	obtained := err.Error()
	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}

	os.Setenv("ProjectID", "ProjectID")
	_, err = getService()
	expected = "ClientSecret must be set"
	obtained = err.Error()
	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
	os.Setenv("ClientSecret", "ClientSecret")
	_, err = getService()
	expected = "google selected but no token provided"
	obtained = err.Error()
	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
	os.Unsetenv("ClientId")
	os.Unsetenv("ProjectID")
	os.Unsetenv("ClientSecret")

}

func TestGetServiceResultFileNotFoundEnv(t *testing.T) {
	os.Setenv("CredentialsFile", "test.json")
	_, err := getService()
	expected := "open test.json: no such file or directory"
	obtained := err.Error()
	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
	os.Unsetenv("CredentialsFile")
}

func TestGetServiceResultFileFoundEnv(t *testing.T) {
	os.Setenv("CredentialsFile", "../tests/credentials.json")
	_, err := getService()
	expected := "google selected but no token provided"
	obtained := err.Error()
	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
	os.Unsetenv("CredentialsFile")
}

func TestGetServiceResultFileFoundBrokenEnv(t *testing.T) {
	os.Setenv("CredentialsFile", "../tests/credentials_broken.json")
	_, err := getService()
	expected := "invalid character '}' after array element"
	obtained := err.Error()
	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
	os.Unsetenv("CredentialsFile")
}

func TestGetGoogleClient(t *testing.T) {
	os.Setenv("google", "true")

	os.Setenv("ClientId", "ClientId")
	os.Setenv("ProjectID", "ProjectID")
	os.Setenv("ClientSecret", "ClientSecret")

	os.Setenv("AccessToken", "true")
	os.Setenv("TokenType", "true")
	os.Setenv("RefreshToken", "true")
	os.Setenv("expiry", "2006-01-02T15:04:05.999999999Z")

	_, err := NewClient()
	obtained := err
	if nil != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", "nil", obtained)
	}
	os.Unsetenv("google")
	os.Unsetenv("ClientId")
	os.Unsetenv("ProjectID")
	os.Unsetenv("ClientSecret")

	os.Unsetenv("AccessToken")
	os.Unsetenv("TokenType")
	os.Unsetenv("RefreshToken")
	os.Unsetenv("expiry")
}

func TestGetServiceResultTokenFileFound(t *testing.T) {
	os.Setenv("TokenFile", "../tests/token.json")
	os.Setenv("CredentialsFile", "../tests/credentials.json")
	_, err := getService()
	expected := "nil"
	obtained := err
	if nil != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
	os.Unsetenv("CredentialsFile")
	os.Unsetenv("TokenFile")
}

func TestGSaveToken(t *testing.T) {
	token, err := tokenFromFile("../tests/token.json")
	if err != nil {
		t.Errorf("\n...expected = %v\n...obtained = %v", "Token to be saved", err.Error())
	}

	err = saveToken("../tests/saveToken.json", token)

	expected := "nil"
	obtained := err

	if nil != err {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
}

func TestGSaveTokenError(t *testing.T) {
	token, err := tokenFromFile("../tests/token.json")
	if err != nil {
		t.Errorf("\n...expected = %v\n...obtained = %v", "Token to be saved", err.Error())
	}

	err = saveToken("../fakepath/saveToken.json", token)

	expected := "open ../fakepath/saveToken.json: no such file or directory"
	obtained := err.Error()

	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
}

func TestGetGoogleClientMissingTokenEnvs(t *testing.T) {
	os.Setenv("google", "true")

	os.Setenv("ClientId", "ClientId")
	os.Setenv("ProjectID", "ProjectID")
	os.Setenv("ClientSecret", "ClientSecret")

	os.Setenv("AccessToken", "true")

	_, err := NewClient()
	obtained := err.Error()
	expected := "TokenType must be set"
	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", "nil", obtained)
	}

	os.Setenv("TokenType", "true")
	_, err = NewClient()
	obtained = err.Error()
	expected = "RefreshToken must be set"
	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", "nil", obtained)
	}

	os.Setenv("RefreshToken", "true")
	_, err = NewClient()
	obtained = err.Error()
	expected = "expiry must be set"
	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", "nil", obtained)
	}

	os.Setenv("expiry", "2006-01-02T15:04:05.999999999Z")
	_, err = NewClient()
	if nil != err {
		t.Errorf("\n...expected = %v\n...obtained = %v", "nil", obtained)
	}

	os.Setenv("expiry", "2006-02T15:04:05.999999999Z")
	_, err = NewClient()
	if nil == err {
		t.Errorf("\n...expected = %v\n...obtained = %v", "Error with parsing date", obtained)
	}

	os.Unsetenv("google")
	os.Unsetenv("ClientId")
	os.Unsetenv("ProjectID")
	os.Unsetenv("ClientSecret")

	os.Unsetenv("AccessToken")
	os.Unsetenv("TokenType")
	os.Unsetenv("RefreshToken")
	os.Unsetenv("expiry")
}
