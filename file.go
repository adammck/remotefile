package remotefile

import (
	"fmt"
	"github.com/adammck/remotefile/iface"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type File struct {
	backend   iface.Backend
	Directory string
}

func New(backend iface.Backend) *File {
	return &File{
		backend:   backend,
		Directory: temporaryDirectory(),
	}
}

func (r *File) Get() (bool, error) {

	// Create the temporary directory to download the file into. Even if the
	// download fails, this must exist for the local file to be written into.

	err := os.MkdirAll(r.Directory, 0700)
	if err != nil {
		return false, err
	}

	exists, rr, err := r.backend.Get()
	if !exists || err != nil {
		return exists, err
	}

	// Download the contents of the remote file into the local file.

	f, err := os.Create(r.Path())
	defer f.Close()
	if err != nil {
		return false, err
	}

	_, err = io.Copy(f, rr)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Put uploads the temporary file to the remote.
func (r *File) Put() error {

	// If the file was deleted, remove it from the remote.
	// TODO: What if the remote already doesn't exist?

	if !pathExists(r.Path()) {
		return r.backend.Delete()
	}

	// Upload the contents of the local file to the remote.

	f, err := os.Open(r.Path())
	if err != nil {
		return err
	}

	return r.backend.Put(f)
}

// Path returns the path (on the local filesystem) of the temporary file which
// the remote file was or will be downloaded to.
func (r *File) Path() string {
	return filepath.Join(r.Directory, r.backend.Filename())
}

// Close deletes the temporary files and directories created by Get.
func (r *File) Close() error {
	return os.RemoveAll(r.Directory)
}

func temporaryDirectory() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("%d%d", time.Now().UnixNano(), r.Int()))
}

// pathExists returns true if the given path exists.
func pathExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}
