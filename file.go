package remotefile

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/adammck/remotefile/iface"
	"github.com/blang/vfs"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type File struct {
	backend   iface.Backend
	Directory string
	fs        vfs.Filesystem
}

func New(backend iface.Backend) *File {
	return &File{
		backend:   backend,
		Directory: temporaryDirectory(),
		fs:        vfs.OS(),
	}
}

func (r *File) Get() (exists bool, err error) {

	// Create the temporary directory to download the file into. Even if the
	// download fails, this must exist for the local file to be written into.

	err = vfs.MkdirAll(r.fs, r.Directory, 0700)
	if err != nil {
		return
	}

	exists, rr, err := r.backend.Get()
	if !exists || err != nil {
		return
	}

	// Download the contents of the remote file into the local file.

	f, err := r.fs.OpenFile(r.Path(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = io.Copy(f, rr)
	if err != nil {
		return
	}

	return
}

// Put uploads the temporary file to the remote.
func (r *File) Put() error {

	// If the file was deleted, remove it from the remote.
	// TODO: What if the remote already doesn't exist?

	if !r.pathExists(r.Path()) {
		return r.backend.Delete()
	}

	// Upload the contents of the local file to the remote.

	f, err := r.fs.OpenFile(r.Path(), os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	return r.backend.Put(f)
}

// Path returns the path (on the local filesystem) of the temporary file which
// the remote file was or will be downloaded to.
func (r *File) Path() string {
	return filepath.Join(r.Directory, r.backend.Filename())
}

// Close deletes the temporary files and directories created by Get.
func (r *File) Close() error {
	return vfs.RemoveAll(r.fs, r.Directory)
}

func temporaryDirectory() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("%d%d", time.Now().UnixNano(), r.Int()))
}

// pathExists returns true if the given path exists.
func (r *File) pathExists(p string) bool {
	_, err := r.fs.Stat(p)
	return err == nil
}

// Checksum returns the hex-encoded SHA1 checksum of the contents of the
// temporary file. This is used to determine whether the file has changed.
// Returns an empty string but no error if the file doesn't exist.
func (r *File) Checksum() (string, error) {
	if !r.pathExists(r.Path()) {
		return "", nil
	}

	f, err := r.fs.OpenFile(r.Path(), os.O_RDONLY, 0600)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha1.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
