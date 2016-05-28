# Remote File

This little Go library provides a common workflow:

1. Download a file from S3
2. Run some program using it
3. Upload the file (if it changed)
4. Clean up

It's similar to Terraform's [remote state][rs], for arbitrary files.


## License

MIT.


[rs]: https://www.terraform.io/docs/state/remote/index.html
