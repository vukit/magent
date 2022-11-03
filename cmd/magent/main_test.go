package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	os.Args = append(os.Args, "-c=../../configs/global.json")
	main()
}
