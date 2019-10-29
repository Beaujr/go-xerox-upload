package client

import "net/http"

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
