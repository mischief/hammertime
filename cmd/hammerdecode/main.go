package main

import (
	"github.com/mischief/hammertime"
	"io"
	"os"
)

func main() {
	dec := hammertime.NewDecoder(os.Stdout)

	io.Copy(dec, os.Stdin)
}
