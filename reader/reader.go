package Reader

import (
	"errors"
	"io"
	"iter"
)

type Reader struct {
	buf                      []byte
	r, w                     int //read and write bytes.
	reader                   io.Reader
	maxEmptyConsecutiveReads int
	err                      error
	readBytes                int
}

const (
	DefaultBufSize                  int = 1024
	DefaultMaxConsecutiveEmptyReads int = 100
	DefaultReadBytes                int = -1
)

var (
	ReachedMaxBytesRead error = errors.New("Reach max bytes read")
)

// readBytes = -1 for no limit.
func NewReaderConfigured(r io.Reader, bufferSize, maxEmptyConsecutiveReads int, readBytes int) *Reader {
	return &Reader{
		buf:                      make([]byte, bufferSize),
		reader:                   r,
		maxEmptyConsecutiveReads: maxEmptyConsecutiveReads,
		readBytes:                readBytes,
	}
}

func NewReader(r io.Reader) *Reader {
	return NewReaderConfigured(r, DefaultBufSize, DefaultMaxConsecutiveEmptyReads, DefaultReadBytes)
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
	rd.readBytes += rd.w
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
		bytesRead := 0
		for {
			if iter.rd.w == iter.rd.r {
				iter.rd.fill()
			}
			if iter.rd.w == iter.rd.r {
				iter.Err = iter.rd.err
				break
			}

			if !yield(iter.rd.buf[iter.rd.r]) {
				break
			}
			iter.rd.r++

			bytesRead++
			if bytesRead >= iter.rd.readBytes {
				iter.Err = ReachedMaxBytesRead
				break
			}
		}
	}
}
