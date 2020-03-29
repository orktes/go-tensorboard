package events

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"io"

	"github.com/pkg/errors"

	"github.com/Applifier/go-tensorflow/types/tensorflow/core/util"
)

const (
	maskDelta = 0xa282ead8

	headerSize = 12
	footerSize = 4
)

var crc32c = crc32.MakeTable(crc32.Castagnoli)

// ErrInvalidChecksum is returned if header or footer checksum is invalid
var ErrInvalidChecksum = errors.New("invalid crc")

// Reader reads TFEvents for an io.Reader
type Reader struct {
	r   io.Reader
	buf *bytes.Buffer
}

// NewReader returns a new reader for given io.Reader
func NewReader(r io.Reader) *Reader {
	return &Reader{
		r:   r,
		buf: bytes.NewBuffer(nil),
	}
}

// Next returns the next event in the reader
func (r *Reader) Next() (*util.Event, error) {
	f := r.r
	buf := r.buf
	buf.Reset()

	_, err := io.CopyN(buf, f, headerSize)
	if err != nil {
		return nil, err
	}

	header := buf.Bytes()

	crc := binary.LittleEndian.Uint32(header[8:12])
	if !verifyChecksum(header[0:8], crc) {
		return nil, errors.Wrap(ErrInvalidChecksum, "length")
	}

	length := binary.LittleEndian.Uint64(header[0:8])
	buf.Reset()

	if _, err = io.CopyN(buf, f, int64(length)); err != nil {
		return nil, err
	}

	if _, err = io.CopyN(buf, f, footerSize); err != nil {
		return nil, err
	}

	payload := buf.Bytes()

	footer := payload[length:]
	crc = binary.LittleEndian.Uint32(footer)
	if !verifyChecksum(payload[:length], crc) {
		return nil, errors.Wrap(ErrInvalidChecksum, "payload")
	}

	ev := &util.Event{}

	return ev, ev.Unmarshal(payload[0:length])
}

func verifyChecksum(data []byte, crcMasked uint32) bool {
	rot := crcMasked - maskDelta
	unmaskedCrc := ((rot >> 17) | (rot << 15))

	crc := crc32.Checksum(data, crc32c)

	return crc == unmaskedCrc
}
