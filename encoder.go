// TODO: check for errors and stuff
package hammertime

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"sync"
	"time"
)

var (
	crap []byte
)

// setup crap slice
func init() {
	crap = make([]byte, 127)
	for i := range crap {
		crap[i] = 0xFF
	}
}

// The hammertime encoder will either send min(buf, 127) bytes as written
// via Write, or send 127 bytes of chaff.
type Encoder struct {
	wr     io.Writer
	wg     sync.WaitGroup
	in     chan []byte
	out    chan []byte
	tick   <-chan time.Time
	errors chan error
}

// Create a new hammertime encoder, writing data or chaff to the
// wr io.Writer.
func NewEncoder(wr io.Writer, tick <-chan time.Time) *Encoder {
	enc := &Encoder{
		wr:     wr,
		in:     make(chan []byte, 1),
		out:    make(chan []byte, 1),
		tick:   tick,
		errors: make(chan error, 1),
	}

	enc.wg.Add(1)

	go func() {
		defer enc.wg.Done()
	loop:
		for {
			select {
			case buf, ok := <-enc.out:
				if !ok {
					break loop
				}
				if _, err := enc.wr.Write(buf); err != nil {
					// TODO: handle error better
					enc.errors <- err
				}
			}
		}
	}()

	enc.wg.Add(1)
	go enc.run()

	return enc
}

func (enc *Encoder) Write(p []byte) (n int, err error) {
	enc.in <- p
	return len(p), nil
}

func (enc *Encoder) Close() error {
	close(enc.in)
	enc.wg.Wait()
	return nil
}

func (enc *Encoder) run() {
	var data []byte

	defer enc.wg.Done()

loop:
	for {
		<-enc.tick
		select {
		case buf, ok := <-enc.in:
			if !ok {
				break loop
			}
			data = append(data, buf...)
		default:
			data = append(data, crap...)
		}

		for len(data) > 0 {
			encoded := new(bytes.Buffer)
			framesize := uint8(math.Min(float64(len(data)), float64(127)))
			binary.Write(encoded, binary.LittleEndian, framesize)
			encoded.Write(data[:framesize])
			data = data[framesize:]
			enc.out <- encoded.Bytes()
		}
	}
	close(enc.out)
}
