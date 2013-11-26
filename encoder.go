// TODO: check for errors and stuff
package hammertime

import (
	"encoding/binary"
	"io"
	"sync"
	"time"
)

// Hammertime encoder. Always sends 127 bytes of data, and will pad
// short writes (< 127 bytes) with chaff.
type Encoder struct {
	wr     io.Writer
	wg     sync.WaitGroup
	in     chan frame
	out    chan frame
	tick   <-chan time.Time
	errors chan error
}

// Create a new hammertime encoder, writing data or chaff to the
// wr io.Writer.
func NewEncoder(wr io.Writer, tick <-chan time.Time) *Encoder {
	enc := &Encoder{
		wr:     wr,
		in:     make(chan frame, 1),
		out:    make(chan frame, 1),
		tick:   tick,
		errors: make(chan error, 1),
	}

	enc.wg.Add(1)

	go func() {
		defer enc.wg.Done()
	loop:
		for {
			buf, ok := <-enc.out
			if !ok {
				break loop
			}

			// TODO: check errors while writing
			binary.Write(enc.wr, binary.LittleEndian, buf.length)
			enc.wr.Write(buf.data)
		}
	}()

	enc.wg.Add(1)
	go enc.encode()

	return enc
}

// Write some data, padded with chaff if len(p) is not divisible by 127.
func (enc *Encoder) Write(p []byte) (n int, err error) {
	for _, f := range makeframes(p) {
		enc.in <- f
	}

	return len(p), nil
}

// Close the encoder. Cleans up any goroutines spawned and closes
// the underlying io.Writer if it is a io.Closer.
func (enc *Encoder) Close() error {
	close(enc.in)
	enc.wg.Wait()
	if closer, ok := enc.wr.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// Encode will run asynchronously to whoever is writing to us,
// adding chaff whenever nothing is on the in channel.
func (enc *Encoder) encode() {
	var fr *frame

	defer enc.wg.Done()

loop:
	for {
		<-enc.tick
		select {
		case buf, ok := <-enc.in:
			if !ok {
				break loop
			}

			fr = &buf
		default:
			fr = &chaff
		}

		enc.out <- *fr

	}
	close(enc.out)
}
