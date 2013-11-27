package hammertime

import (
	"bytes"
	"testing"
	"time"
)

func TestDecoder(t *testing.T) {
	testdata := []byte("hammertime")
	output := new(bytes.Buffer)
	enc := NewEncoder(NewDecoder(output), time.Tick(10*time.Millisecond))

	enc.Write(testdata)
	time.Sleep(50 * time.Millisecond)
	enc.Write(testdata)
	enc.Close()

	t.Logf("output string: %s", output.String())
	t.Logf("output: %X", output.Bytes())
	if bytes.Contains(output.Bytes(), chaff.data) {
		t.Error("decoded output has chaff")
	}

	if !bytes.Contains(output.Bytes(), testdata) {
		t.Error("doesn't have testdata")
	}
}
