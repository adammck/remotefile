package iface

import "io"

type Backend interface {

	// Get returns a bool indicating whether the remote file exists, and reader
	// to read its contents.
	Get() (bool, io.Reader, error)

	// Put writes the contents of a ReadSeeker to the remote file.
	Put(io.ReadSeeker) error

	// Delete removes the remote file.
	Delete() error

	// Filename returns the filename (the last path element) which should be
	// used if the remote file is written to disk locally. This ensures that
	// programs which expect a specific filename or extension aren't confused.
	Filename() string
}
