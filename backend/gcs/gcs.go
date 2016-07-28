package gcs

import (
	"bytes"
	"io"
	"net/http"
	"path"

	"github.com/adammck/remotefile/iface"

	"google.golang.org/api/googleapi"
	storage "google.golang.org/api/storage/v1"
)

type GCS struct {
	service *storage.Service
	Bucket  string
	Path    string
}

var _ iface.Backend = (*GCS)(nil)

//
// See: https://godoc.org/golang.org/x/oauth2/google
//
func New(client *http.Client, bucket string, path string) (*GCS, error) {
	svc, err := storage.New(client)
	if err != nil {
		return nil, err
	}

	return &GCS{
		service: svc,
		Bucket:  bucket,
		Path:    path,
	}, nil
}

func (g *GCS) Get() (bool, io.Reader, error) {
	res, err := g.service.Objects.Get(g.Bucket, g.Path).Download()
	if err != nil {

		// First try to cast the error back to a google Error, so we can extract
		// the http status. We don't consider 404 to be an error.
		gerr, ok := err.(*googleapi.Error)
		if ok && gerr.Code == 404 {
			return false, &bytes.Buffer{}, nil
		}

		// Otherwise, return the error as-is.
		return false, &bytes.Buffer{}, err
	}

	return true, res.Body, nil
}

func (g *GCS) Put(r io.ReadSeeker) error {
	_, err := g.service.Objects.Insert(g.Bucket, &storage.Object{Name: g.Path}).Media(r).Do()
	return err
}

func (g *GCS) Delete() error {
	return g.service.Objects.Delete(g.Bucket, g.Path).Do()
}

func (g *GCS) Filename() string {
	return path.Base(g.Path)
}
