package hammertime

import (
	"io"
)

type decState int

const (
	decHdr decState = iota
	decFrame
)

// Hammertime decoder.
type Decoder struct {
	wr io.Writer

	// current frame being decoded
	compliance frame

	// how many bytes of the frame we have
	gotlen uint8

	// state, getting header or gettin frame content
	state decState
}

func NewDecoder(wr io.Writer) *Decoder {
	dec := &Decoder{
		wr:         wr,
		compliance: frame{data: make([]byte, maxframe)},
	}

	return dec
}

// Strip chaff from a hammertime stream.
func (dec *Decoder) Write(p []byte) (n int, err error) {
	for _, b := range p {
		switch dec.state {
		case decHdr:
			dec.compliance.length = b & 127
			dec.gotlen = 0
			dec.state = decFrame
		case decFrame:
			if int(dec.gotlen) < cap(dec.compliance.data) {
				if dec.gotlen < dec.compliance.length {
					dec.compliance.data[dec.gotlen] = b
				}
				dec.gotlen++
			}

			if dec.gotlen == dec.compliance.length {
				// we have a whole frame
				if _, err := dec.wr.Write(dec.compliance.data[:dec.compliance.length]); err != nil {
					return 0, err
				}
			}

			if int(dec.gotlen) == cap(dec.compliance.data) {
				dec.state = decHdr
			}

		}
	}

	return len(p), err

	// old stuff
	/*
		for len(p) > 0 {
			switch dec.state {
			case DecHdr:
				fmt.Fprintf(os.Stderr, "p: (%d) %X\n", len(p), p)
				framehdr := p[0]
				dec.compliance.length = p[0]
				// remove when done debugging

				frametype := framehdr & 128
				framelen := framehdr & 127
				fmt.Fprintf(os.Stderr, "frame type %d len %d\n", frametype, framelen)

				p = p[1:]
				dec.gotlen = copy(dec.compliance.data, p)
				dec.state = DecFrame
			case DecFrame:
				p = p[:copy(dec.compliance.data[len(dec.compliance.data):], p)]

				fmt.Fprintf(os.Stderr, "p: (%d) %X\n", len(p), p)
				fmt.Fprintf(os.Stderr, "length: %d\n", dec.compliance.length&127)
				fmt.Fprintf(os.Stderr, "data: (%d) %X\n", len(dec.compliance.data), dec.compliance.data)

				if dec.compliance.length&128 == 0 && len(dec.compliance.data) == cap(dec.compliance.data) {
					_, err = dec.wr.Write(dec.compliance.data[:dec.compliance.length&127])
					dec.state = DecHdr
				}
			}

		}
	*/
}

// Close the decoder, and the underlying writer if it is an
// io.Closer.
func (dec *Decoder) Close() error {
	if closer, ok := dec.wr.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
