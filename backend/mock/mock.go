package mock

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/adammck/remotefile/iface"
)

type Mock struct {
	fn         string
	RemoteData []byte
}

func New(fn string) *Mock {
	return &Mock{
		fn:         fn,
		RemoteData: nil,
	}
}

func (m *Mock) Get() (bool, io.Reader, error) {
	if m.RemoteData == nil {
		return false, bytes.NewReader(nil), nil
	}

	return true, bytes.NewReader(m.RemoteData), nil
}

func (m *Mock) Put(r io.ReadSeeker) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	m.RemoteData = b
	return nil
}

func (m *Mock) Delete() error {
	m.RemoteData = nil
	return nil
}

func (m *Mock) Filename() string {
	return m.fn
}

// Ensure that Mock implements the Backend interface.
var _ iface.Backend = (*Mock)(nil)
