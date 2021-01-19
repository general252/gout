package uarray

import (
	"bytes"
	"fmt"
	"io"
)

// from https://github.com/orcaman/writerseeker
// add method Insert(p []byte)

// WriterSeeker is an in-memory io.WriteSeeker implementation
type WriterSeeker struct {
	buf bytes.Buffer
	pos int
}

var _ io.WriteSeeker = new(WriterSeeker)

// Write writes to the buffer of this WriterSeeker instance
func (ws *WriterSeeker) Write(p []byte) (n int, err error) {
	// If the offset is past the end of the buffer, grow the buffer with null bytes.
	if extra := ws.pos - ws.buf.Len(); extra > 0 {
		if _, err := ws.buf.Write(make([]byte, extra)); err != nil {
			return n, err
		}
	}

	// If the offset isn't at the end of the buffer, write as much as we can.
	if ws.pos < ws.buf.Len() {
		n = copy(ws.buf.Bytes()[ws.pos:], p)
		p = p[n:]
	}

	// If there are remaining bytes, append them to the buffer.
	if len(p) > 0 {
		var bn int
		bn, err = ws.buf.Write(p)
		n += bn
	}

	ws.pos += n
	return n, err
}

// Insert insert p at pos
func (ws *WriterSeeker) Insert(p []byte) (n int, err error) {
	var buffLen = ws.buf.Len()

	if extra := ws.pos - ws.buf.Len(); extra > 0 {
		if _, err := ws.buf.Write(make([]byte, extra)); err != nil {
			return n, err
		}
	}
	if _, err := ws.buf.Write(make([]byte, len(p))); err != nil {
		return n, err
	}

	if ws.pos < ws.buf.Len() {
		var buf = ws.buf.Bytes()

		var tmpAfterPosBuf []byte
		if ws.pos < buffLen {
			tmpAfterPosBuf = make([]byte, buffLen-ws.pos)
			copy(tmpAfterPosBuf, buf[ws.pos:buffLen])
		}

		n = copy(buf[ws.pos:], p)

		if tmpAfterPosBuf != nil {
			copy(buf[ws.pos+len(p):], tmpAfterPosBuf)
		}
	}

	ws.pos += n
	return n, err
}

// Seek seeks in the buffer of this WriterSeeker instance
func (ws *WriterSeeker) Seek(offset int64, whence int) (int64, error) {
	newPos, offs := 0, int(offset)
	switch whence {
	case io.SeekStart:
		newPos = offs
	case io.SeekCurrent:
		newPos = ws.pos + offs
	case io.SeekEnd:
		newPos = ws.buf.Len() + offs
	}
	if newPos < 0 {
		return 0, fmt.Errorf("negative result pos")
	}
	ws.pos = newPos
	return int64(newPos), nil
}

// Close :
func (ws *WriterSeeker) Close() error {
	return nil
}

func (ws *WriterSeeker) Bytes() []byte {
	return ws.buf.Bytes()
}
func (ws *WriterSeeker) String() string {
	return ws.buf.String()
}
