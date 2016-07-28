# Remote File
[![GoDoc](https://godoc.org/github.com/adammck/remotefile?status.svg)](https://godoc.org/github.com/adammck/remotefile)
[![Build Status](https://travis-ci.org/adammck/remotefile.svg?branch=master)](https://travis-ci.org/adammck/remotefile)

This little Go library provides a common workflow:

1. Download a file from S3
2. Run some program using it
3. Upload the file (if it changed)
4. Clean up

It's similar to Terraform's [remote state][rs], for arbitrary files.


## Usage

```go
svc := s3.New(session.New())
backend := s3be.New(svc, "my-bucket", "path/to/file.txt")
rf := remotefile.New(backend)

// download the file to tmp
exists, err := rf.Get()
if err != nil {
  exit("error pulling file: %s", err)
}

// edit the file locally
cmd := exec.Command("vim", rf.Path())
cmd.Stdin = os.Stdin
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
err = cmd.Run()
if err != nil {
  exit("vim returned an error: %s", err)
}

// upload the file
err = f.Put()
if err != nil {
  exit("error pushing file: %s", err)
}

// delete the file (from the local filesystem)
err = f.Close()
if err != nil {
  exit("error cleaning up: %s", err)
}
```

## Testing

To run integration tests:

```bash
export INTEGRATION_TESTS=1
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
export GOOGLE_PROJECT=
go test -v ./backend/...
```

These aren't run by Travis, so please be sure to run them before merging changes
to master. They also cost (a miniscule amount of) money, because they create and
modify remote resources, so be sure to clean up manually when they fail.

## License

MIT.


[rs]: https://www.terraform.io/docs/state/remote/index.html
