package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/adammck/remotefile"
	s3be "github.com/adammck/remotefile/backend/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	svc := s3.New(session.New())
	f := remotefile.New(s3be.New(svc, "adammck", "remotefile/edit.txt"))

	fmt.Println("downloading...")
	_, err := f.Get()
	checkErr(err)

	p := f.Path()
	fmt.Printf("editing %s\n", p)
	cmd := exec.Command("vim", p)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	fmt.Println("uploading...")
	err = f.Put()
	checkErr(err)

	fmt.Println("cleaning up...")
	err = f.Close()
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
