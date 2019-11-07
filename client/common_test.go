package client

import (
	"os"
	"testing"
)

func TestGetEmptyEnv(t *testing.T) {
	os.Setenv("CredentialsFile", "")
	_, err := getEnvVar("CredentialsFile")
	expected := "CredentialsFile must not be empty"
	obtained := err.Error()
	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
	os.Unsetenv("CredentialsFile")

}

func TestGetFileSystemClient(t *testing.T) {
	os.Setenv("PGID", "1000")
	os.Setenv("GID", "1000")
	_, err := NewClient()
	obtained := err
	if nil != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", "nil", obtained)
	}
	os.Unsetenv("PGID")
	os.Unsetenv("GID")

}
func TestGetFileSystemClientNoGID(t *testing.T) {
	os.Setenv("PGID", "1000")
	_, err := NewClient()
	expected := "GID must be set"
	obtained := err.Error()
	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
	os.Unsetenv("PGID")

}
func TestGetFileSystemClientNoPGID(t *testing.T) {
	os.Setenv("GID", "1000")
	_, err := NewClient()
	expected := "PGID must be set"
	obtained := err.Error()
	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
	os.Unsetenv("GID")

}
