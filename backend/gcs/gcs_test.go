package gcs

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/mattetti/filebuffer"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

func TestEverything(t *testing.T) {

	// This test creates a bucket and file(s) on GCP, which costs actual money
	// (albeit not a lot of it), so only run if explicitly enabled.

	if os.Getenv("INTEGRATION_TESTS") == "" {
		t.Skipf("skipping; INTEGRATION_TESTS must be set")
	}

	// Verify that the environment set up correctly.

	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		t.Skipf("skipping; GOOGLE_APPLICATION_CREDENTIALS must be set")
	}

	project := os.Getenv("GOOGLE_PROJECT")
	if project == "" {
		t.Skipf("skipping; GOOGLE_PROJECT must be set")
	}

	t.Logf("running gcs integration tests...")

	// -------------------------------------------------------------------------

	client, err := google.DefaultClient(oauth2.NoContext, storage.DevstorageFullControlScope)
	if err != nil {
		t.Fatalf("error creating http client: %q", err)
	}

	svc, err := storage.New(client)
	if err != nil {
		t.Fatalf("error creating storage service: %q", err)
	}

	bucket := fmt.Sprintf("test-remotefile-%d", time.Now().UnixNano())

	// Create a bucket to write to, and delete it when we're finished.

	_, err = svc.Buckets.Insert(project, &storage.Bucket{Name: bucket}).Do()
	if err != nil {
		t.Fatalf("error when creating bucket: %q", err)
	}
	defer func() {
		err := svc.Buckets.Delete(bucket).Do()
		if err != nil {
			t.Logf("error deleting bucket: %q", err)
		}
	}()

	// We need filebuffer here rather bytes.Buffer, because the latter is not a
	// ReadSeeker.
	one := []byte("one")
	two := []byte("two")

	// -------------------------------------------------------------------------

	be, err := New(client, bucket, "/test.txt")
	if err != nil {
		t.Fatalf("error when instantiating gcs backend: %q", err)
	}

	// 1. get file which doesn't exist
	exists, body, err := be.Get()
	assert.NoError(t, err)
	assert.False(t, exists)
	b, _ := ioutil.ReadAll(body)
	assert.Equal(t, []byte{}, b)

	// 2. put file
	err = be.Put(filebuffer.New(one))
	assert.NoError(t, err)

	// 3. get contents
	exists, body, err = be.Get()
	assert.NoError(t, err)
	assert.True(t, exists)
	b, _ = ioutil.ReadAll(body)
	assert.Equal(t, one, b)

	// 4. overwrite file
	err = be.Put(filebuffer.New(two))
	assert.NoError(t, err)

	// 5. get new contents
	exists, body, err = be.Get()
	assert.NoError(t, err)
	assert.True(t, exists)
	b, _ = ioutil.ReadAll(body)
	assert.Equal(t, two, b)

	// 6. delete file
	err = be.Delete()
	assert.NoError(t, err)
}
