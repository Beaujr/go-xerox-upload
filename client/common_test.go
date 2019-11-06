package client

import (
	"testing"
	"os"
)

func TestGetEmptyEnv(t *testing.T) {
	os.Setenv("CredentialsFile", "")
	_, err := getEnvVar("CredentialsFile")
	expected := "open test.json: no such file or directory"
	obtained := err.Error()
	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
}

