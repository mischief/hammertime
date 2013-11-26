package hammertime

const (
	maxframe = 127
)

type frame struct {
	length uint8
	data   []byte
}

var (
	ChaffByte = byte('~')
	chaff     frame
)

// setup crap frame
func init() {
	chaff = frame{length: 0x0, data: make([]byte, maxframe)}
	for i := range chaff.data {
		// XXX: use random data?
		chaff.data[i] = ChaffByte
	}
}

func makeframes(p []byte) []frame {
	var out []frame

	// if we have no data, we have no frames
	if len(p) <= 0 {
		return nil
	}

	// if we have data larger than maxframe, chop it up
	for len(p) > maxframe {
		f := frame{length: maxframe, data: make([]byte, maxframe)}

		// copies maxframe bytes from p into f.data
		copy(f.data, p)
		out = append(out, f)

		p = p[maxframe:]
	}

	// now if we have anything left over,
	// make a partial frame and fill the end with chaff
	if len(p) > 0 {
		lastframe := frame{length: uint8(len(p)), data: make([]byte, maxframe)}
		copy(lastframe.data, p)
		copy(lastframe.data[len(p):], chaff.data)
		out = append(out, lastframe)
	}

	return out
}
