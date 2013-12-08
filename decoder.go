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
}

// Close the decoder
func (dec *Decoder) Close() error {
	return nil
}
