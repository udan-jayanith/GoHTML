package Reader

import (
	"errors"
	"io"
	"iter"
)

type Reader struct {
	buf                      []byte
	r, w                     int //read and written bytes.
	reader                   io.Reader
	maxEmptyConsecutiveReads int
	err                      error
	maxByteReads, byteReads  int
}

const (
	DefaultBufSize                  int = 1024
	DefaultMaxConsecutiveEmptyReads int = 100
	DefaultByteReads                int = -1
)

var (
	ReachedMaxBytesRead error = errors.New("Reached max bytes read")
)

// maxByteReads = -1 for no limit.
func NewReaderConfigured(r io.Reader, bufferSize, maxEmptyConsecutiveReads int, maxByteReads int) *Reader {
	return &Reader{
		buf:                      make([]byte, bufferSize),
		reader:                   r,
		maxEmptyConsecutiveReads: maxEmptyConsecutiveReads,
		maxByteReads:             maxByteReads,
	}
}

func NewReader(r io.Reader) *Reader {
	return NewReaderConfigured(r, DefaultBufSize, DefaultMaxConsecutiveEmptyReads, DefaultByteReads)
}

func (rd *Reader) fill() {
	if rd.r != rd.w {
		panic("Must first read available all the data before calling fill.")
	}

	count := 1
	for count <= rd.maxEmptyConsecutiveReads {
		n, err := rd.reader.Read(rd.buf)
		rd.w = n
		rd.r = 0
		if err != nil {
			rd.err = err
			break
		} else if n != 0 {
			break
		}
		count++
	}
	if count == rd.maxEmptyConsecutiveReads {
		rd.err = io.EOF
	}
	rd.byteReads += rd.w
}

type Iter struct {
	Err error
	rd  *Reader
}

func NewIter(rd *Reader) Iter {
	return Iter{
		rd: rd,
	}
}

func (iter *Iter) Loop() iter.Seq[byte] {
	return func(yield func(byt byte) bool) {
		for {
			if iter.rd.w == iter.rd.r {
				iter.rd.fill()
			}
			if iter.rd.w == iter.rd.r {
				iter.Err = iter.rd.err
				break
			}

			if !yield(iter.rd.buf[iter.rd.r]) {
				iter.rd.r++
				break
			}
			iter.rd.r++
		}
	}
}
