package hammertime

import (
	"bytes"
	"testing"
)

func TestMakeFrameLong(t *testing.T) {
	longlen := 200
	long := make([]byte, longlen)
	for i := range long {
		long[i] = byte(i)
	}

	frames := makeframes(long)

	t.Logf("long in:  (%d) %X", len(long), long)
	for _, f := range frames {
		t.Logf("long out: (%d) %X", f.length, f.data)
	}

	if len(frames) != 2 {
		t.Errorf("expected 2 frames from %d bytes of data, got %d frames", len(long), len(frames))
	}

	if !bytes.Equal(frames[0].data, long[:maxframe]) {
		t.Errorf("full frame doesn't contain original data!")
	}

	if !bytes.Equal(frames[1].data[:longlen-maxframe], long[maxframe:]) {
		t.Errorf("partial frame doesn't have partial data")
		t.Logf("frame[1]: %X", frames[1].data[:longlen-maxframe])
		t.Logf("long[maxframe:]: %X", long[maxframe:])
	}
}

func TestMakeFrameShort(t *testing.T) {
	shortlen := 100
	short := make([]byte, shortlen)
	for i := range short {
		short[i] = byte(i)
	}

	frames := makeframes(short)

	t.Logf("short in:  (%d) %X", len(short), short)
	for _, f := range frames {
		t.Logf("short out: (%d) %X", f.length, f.data)
	}

	if len(frames) != 1 {
		t.Errorf("expected 1 frame from data of size < maxframe, got %d", len(frames))
	}

	if !bytes.Equal(frames[0].data[:shortlen], short[:shortlen]) {
		t.Errorf("short frame has wrong content")
	}

	cmpchaff := make([]byte, maxframe-shortlen)
	for i := range cmpchaff {
		cmpchaff[i] = ChaffByte
	}

	if !bytes.Equal(frames[0].data[shortlen:], cmpchaff) {
		t.Errorf("end of short frame doesn't have the right chaff")
	}
}
