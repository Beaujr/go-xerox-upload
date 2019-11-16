package client

import "testing"

var fsClient = NewFileSystemClient(501, 20)

func TestFilesystem_CleanPath(t *testing.T) {
	obtained := fsClient.CleanPath("////this//is/a///test\\")

	expected := "/this/is/a/test/"

	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
}

func TestFilesystem_MakeListDelDirectory(t *testing.T) {
	obtained := fsClient.MakeDirectory("../tests/dir")
	expected := "nil"

	if obtained != nil {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}

	_, err := fsClient.ListDirectory("../tests")
	if err != nil {
		t.Errorf("\n...expected = %v\n...obtained = %v", "nil", err.Error())
	}

	_, err = fsClient.ListDirectory("../nofolder")
	if err == nil {
		t.Errorf("\n...expected = %v\n...obtained = %v", "../nofolder: no such file or directory", "nil")
	}

	obtained = fsClient.DeleteDir("../tests/dir")
	expected = "nil"

	if obtained != nil {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
}

func TestFilesystem_DelDirectory(t *testing.T) {
	obtained := fsClient.MakeDirectory("../tests/dir")
	expected := "nil"

	if obtained != nil {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}

	obtained = fsClient.DeleteDir("../tests/dir")
	expected = "nil"

	if obtained != nil {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}

	err := fsClient.DeleteDir("../tests/dir")
	expected = "remove ../tests/dir: no such file or directory"

	if expected != err.Error() {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
}

//func TestFilesystem_MakeDirectoryErrors(t *testing.T) {
//err := fsClient.MakeDirectory("")
//expected := "mkdir : no such file or directory"
//
//
//obtained := err
//if obtained == nil {
//	t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
//}
//
//nfc := NewFileSystemClient(0, 0)
//err = nfc.MakeDirectory("../tests/notreal")
//expected = "chown ../tests/notreal: operation not permitted"
//obtained = err
//
//if obtained == nil {
//	t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
//}
//
//err = nfc.DeleteDir("../tests/notreal")
//if err != nil {
//	t.Errorf("Failure to delete test dir")
//}

//
//obtained = fsClient.DeleteDir("../tests/dir")
//expected = "nil"
//
//if obtained != nil {
//	t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
//}
//
//err := fsClient.DeleteDir("../tests/dir")
//expected = "remove ../tests/dir: no such file or directory"
//
//if expected != err.Error() {
//	t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
//}
//}
