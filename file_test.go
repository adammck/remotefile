package remotefile

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/adammck/remotefile/iface"
	"github.com/blang/vfs"
	"github.com/blang/vfs/memfs"
	"github.com/stretchr/testify/assert"
)

func TestFile_Get(t *testing.T) {
	f, backend, tmpDir, fs := newFile()

	// When remote file doesn't exist...

	exists, err := f.Get()
	assert.NoError(t, err)
	assert.False(t, exists)

	// expect the temporary directory to have been created,

	info, err := fs.Stat(tmpDir)
	if err != nil {
		t.Fatalf("expected no error when statting tmpDir, got %q", err)
	}
	if !info.IsDir() {
		t.Errorf("expected directory at %q, got file with mode %d", tmpDir, info.Mode())
	}

	// but expect no local file to have been created.

	_, err = fs.OpenFile(f.Path(), os.O_RDONLY, 0)
	assert.Error(t, err)

	// When remote file does exist...

	backend.RemoteFile = []byte("hello, world")
	exists, err = f.Get()
	assert.NoError(t, err)
	assert.True(t, exists)

	// expect local file with same contents

	localFile, err := fs.OpenFile(f.Path(), os.O_RDONLY, 0)
	assert.NoError(t, err)
	contents, err := ioutil.ReadAll(localFile)
	if err != nil {
		t.Fatalf("expected no error when reading local file, got %q", err)
	}
	assert.Equal(t, backend.RemoteFile, contents)
}

// ----

func newFile() (*File, *MockBackend, string, vfs.Filesystem) {
	backend := NewMockBackend("test.txt")
	tmpDir := "/tmp/a/b/c"
	fs := memfs.Create()

	return &File{
		backend:   backend,
		Directory: tmpDir,
		fs:        fs,
	}, backend, tmpDir, fs
}

// ----

type MockBackend struct {
	fn         string
	RemoteFile []byte
}

func NewMockBackend(fn string) *MockBackend {
	return &MockBackend{
		fn:         fn,
		RemoteFile: nil,
	}
}

func (m *MockBackend) Get() (bool, io.Reader, error) {
	if m.RemoteFile == nil {
		return false, bytes.NewReader(nil), nil
	}

	return true, bytes.NewReader(m.RemoteFile), nil
}

func (m *MockBackend) Put(r io.ReadSeeker) error {
	panic("MockBackend.Put: not implemented")
}

func (m *MockBackend) Delete() error {
	panic("MockBackend.Delete: not implemented")
}

func (m *MockBackend) Filename() string {
	return m.fn
}

// Ensure that MockBackend implements the Backend interface.
var _ iface.Backend = (*MockBackend)(nil)
