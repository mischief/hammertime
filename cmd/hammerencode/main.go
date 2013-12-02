package main

import (
	"flag"
	"github.com/mischief/hammertime"
	"io"
	"os"
	"time"
)

var ms = flag.Int("ms", 20, "packet frequency in ms")

func main() {
	flag.Parse()
	enc := hammertime.NewEncoder(os.Stdout, time.Tick(time.Duration(*ms)*time.Millisecond))

	io.Copy(enc, os.Stdin)
}
