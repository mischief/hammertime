package hammertime

import (
	"bytes"
	"testing"
	"time"
)

func TestEncoder(t *testing.T) {
	testdata := []byte("hammertime")
	output := new(bytes.Buffer)
	enc := NewEncoder(output, time.Tick(10*time.Millisecond))

	enc.Write(testdata)
	time.Sleep(50 * time.Millisecond)
	enc.Write(testdata)
	enc.Close()

	t.Logf("output: %X", output.Bytes())
	if !bytes.Contains(output.Bytes(), crap) {
		t.Error("doesn't have any crap")
	}

	if !bytes.Contains(output.Bytes(), testdata) {
		t.Error("doesn't have testdata")
	}
}
