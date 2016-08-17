package main

import (
	"fmt"
	"os"
)

func fatalf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func fatalln(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}
