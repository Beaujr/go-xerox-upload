package client

import (
	"testing"
	"os"
)

func TestGetServiceResultError(t *testing.T) {
	//rec := httptest.NewRecorder()
	//req, _ := http.NewRequest(http.MethodPost, "/v1/user/FakeUserId/entitlement", bytes.NewBuffer([]byte(validEntitlementJson)))
	//env := Env{db: &mockSuccessAdminDB{}}
	//
	//test_helpers.InvokeHandler(http.Handler(createEntitlement(&env)), entitlementEndpointPath, rec, req)
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
}

func TestGetServiceResultFileNotFoundEnv(t *testing.T) {
	os.Setenv("CredentialsFile", "test.json")
	_, err := getService()
	expected := "open test.json: no such file or directory"
	obtained := err.Error()
	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
}

func TestGetServiceResultFileFoundEnv(t *testing.T) {
	os.Setenv("CredentialsFile", "../tests/credentials.json")
	_, err := getService()
	expected := "google selected but no token provided"
	obtained := err.Error()
	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
}