package main

import (
	"fmt"
	"os"
)

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage: just-s3 cp [source] [destination]")
	fmt.Fprintln(os.Stderr, " the source or destination may be an s3 URL of")
	fmt.Fprintln(os.Stderr, " the form s3://bucket/objectname")
}

func main() {
	if len(os.Args) < 3 || os.Args[1] != "cp" {
		printUsage()
		return
	}
	if err := copyObject(os.Args[2], os.Args[3], NewAwsFactory()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
