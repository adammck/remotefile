# Remote File

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

## License

MIT.


[rs]: https://www.terraform.io/docs/state/remote/index.html
