package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"

	"github.com/adammck/remotefile"
	s3be "github.com/adammck/remotefile/backend/s3"
	"github.com/adammck/remotefile/iface"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("usage: %s [url]\n", os.Args[0])
		os.Exit(1)
	}

	u, err := url.Parse(os.Args[1])
	checkErr(err)

	var be iface.Backend
	switch u.Scheme {
	case "s3":
		// The S3 backend is configured via ENV
		svc := s3.New(session.New())
		be = s3be.New(svc, u.Host, u.Path)
	default:
		checkErr(fmt.Errorf("invalid scheme: %s", u.Scheme))
	}

	f := remotefile.New(be)

	// Delete temp files when we're done, whatever happens.
	defer func() {
		err := f.Close()
		checkErr(err)
	}()

	// Download the file, if it exists on the remote.
	_, err = f.Get()
	checkErr(err)

	// Store the current state, to compare after edit.
	chk1, err := f.Checksum()
	checkErr(err)

	// Edit the temp file.
	p := f.Path()
	cmd := exec.Command("vim", p)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	// Stop if the file wasn't changed.
	chk2, err := f.Checksum()
	checkErr(err)
	if chk2 == chk1 {
		return
	}

	// Upload the temp file to replace the remote, or delete it if the temp file
	// was deleted.
	err = f.Put()
	checkErr(err)
}

// checkErr exits the program if err is not nil.
func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
