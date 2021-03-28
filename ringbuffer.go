package main

import (
	"errors"
)

type RingBuffer struct {
	buf  []byte
	pos  int
	size int
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		buf:  make([]byte, size),
		pos:  0,
		size: size,
	}

}
func (r *RingBuffer) Write(p []byte) (n int, err error) {

	for _, b := range p {
		r.buf[r.pos] = b
		r.pos = (r.pos + 1) % r.size
	}
	return len(p), nil
}

func (r *RingBuffer) LookBack(dist, len int) ([]byte, error) {
	if dist > r.size {
		return nil, errors.New("dist larger than buffersize")
	}
	out := make([]byte, len)

	curPos := (r.pos - dist) % r.size
	if curPos < 0 {
		curPos += r.size
	}
	for i := 0; i < len; i++ {
		out[i] = r.buf[curPos]
		curPos = (curPos + 1) % r.size
	}
	return out, nil
}
