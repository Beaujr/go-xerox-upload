package client

import "net/http"

// ListDirectory is the Payload value from the Printer to List Directory Values to avoid filename collisions
const ListDirectory = "ListDir"

// MakeDir is the Payload value from the Printer to Make a Directory if it doesn't exit
const MakeDir = "MakeDir"

// PutFile is the Payload value from the Printer to upload a file
const PutFile = "PutFile"

// DeleteFile is the Payload value from the Printer to Delete a file normally a *LCK file
const DeleteFile = "DeleteFile"

// RemoveDir is the Payload value from the Printer to Delete a directory
const RemoveDir = "RemoveDir"

// DestDir is the Payload value from the Printer for the files dir on the filesystem
const DestDir = "destDir"

// DestName is the Payload value from the Printer for the filename  on the filesystem
const DestName = "destName"

// Operation is the Payload field from the Printer for the the operation to happen, ListDir, MakeDir etc
const Operation = "theOperation"

// Sendfile is the Payload field from the Printer for the the file itself
const Sendfile = "sendfile"

// XRXNOTFOUND is the not found error message
const XRXNOTFOUND = "XRXNOTFOUND"

// XRXERROR is the default error message
const XRXERROR = "XRXERROR"

// XRXDIREXISTS is the directory exists already
const XRXDIREXISTS = "XRXDIREXISTS"

// XRXBADNAME is the filename is bad due to to FS constraints
const XRXBADNAME = "XRXBADNAME"

// XeroxApi Interface for all Printer to Server Interactions
type XeroxApi interface {
	ListDirectory(directory string) (string, error)
	CleanPath(directory string) string
	DeleteDir(directory string) error
	PutFile(r *http.Request, directory string) ([]byte, error)
	MakeDirectory(directory string) error
}
