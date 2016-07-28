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

func TestFileGet(t *testing.T) {
	f, backend, tmpDir, fs := newFile(nil)

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

	backend.RemoteData = []byte("hello, world")
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
	assert.Equal(t, backend.RemoteData, contents)
}

func TestFilePut(t *testing.T) {
	f, backend, _, _ := newFile([]byte("remote data"))

	// When local file doesn't exist, the remote should be deleted.

	err := f.Put()
	assert.NoError(t, err)
	assert.Nil(t, backend.RemoteData)

	// When local file exists, remote should be updated.

	f, backend, _, fs := newFile(nil)
	localData := []byte("hello, world")

	err = vfs.MkdirAll(fs, f.Directory, 0700)
	if err != nil {
		t.Fatalf("expected no error when creating tmp dir, got %q", err)
	}

	err = writeFile(fs, f.Path(), localData, 0600)
	if err != nil {
		t.Fatalf("expected no error when writing local file, got %q", err)
	}

	err = f.Put()
	assert.NoError(t, err)
	assert.Equal(t, localData, backend.RemoteData)
}

func TestFileChecksum(t *testing.T) {
	f, _, _, fs := newFile(nil)
	local := []byte("pup in a cup")

	err := vfs.MkdirAll(fs, f.Directory, 0700)
	if err != nil {
		t.Fatalf("expected no error when creating tmp dir, got %q", err)
	}

	sum, err := f.Checksum()
	assert.NoError(t, err)
	assert.Equal(t, "", sum)

	err = writeFile(fs, f.Path(), local, 0600)
	if err != nil {
		t.Fatalf("expected no error when writing local file, got %q", err)
	}

	sum, err = f.Checksum()
	assert.NoError(t, err)
	assert.Equal(t, "2e3b1d2f993e7df69e9fb761f0b9434bfec2e44c", sum)
}

// ----

func newFile(remoteData []byte) (*File, *MockBackend, string, vfs.Filesystem) {
	backend := NewMockBackend("test.txt")
	tmpDir := "/tmp/a/b/c"
	fs := memfs.Create()

	return &File{
		backend:   backend,
		Directory: tmpDir,
		fs:        fs,
	}, backend, tmpDir, fs
}

// port of ioutil.Writefile for vfs
func writeFile(fs vfs.Filesystem, filename string, data []byte, perm os.FileMode) error {
	f, err := fs.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// ----

type MockBackend struct {
	fn         string
	RemoteData []byte
}

func NewMockBackend(fn string) *MockBackend {
	return &MockBackend{
		fn:         fn,
		RemoteData: nil,
	}
}

func (m *MockBackend) Get() (bool, io.Reader, error) {
	if m.RemoteData == nil {
		return false, bytes.NewReader(nil), nil
	}

	return true, bytes.NewReader(m.RemoteData), nil
}

func (m *MockBackend) Put(r io.ReadSeeker) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	m.RemoteData = b
	return nil
}

func (m *MockBackend) Delete() error {
	m.RemoteData = nil
	return nil
}

func (m *MockBackend) Filename() string {
	return m.fn
}

// Ensure that MockBackend implements the Backend interface.
var _ iface.Backend = (*MockBackend)(nil)
