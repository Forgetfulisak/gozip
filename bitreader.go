package main

import (
	"bufio"
	"errors"
	"log"
)

// Stolen from here: https://play.golang.org/p/Wyr_K9YAro
type BitReader struct {
	reader *bufio.Reader
	byte   byte
	offset byte
}

func NewBitReader(r *bufio.Reader) *BitReader {
	return &BitReader{r, 0, 0}
}

func (r *BitReader) ReadBit() (bool, error) {
	if r.offset == 8 {
		r.offset = 0
	}
	if r.offset == 0 {
		var err error
		// fmt.Println("reading")
		if r.byte, err = r.reader.ReadByte(); err != nil {
			log.Fatalln("error!")
			return false, err
		}
	}
	//bit := (r.byte & (0x80 >> r.offset)) != 0
	bit := ((r.byte >> r.offset) & 0x1) != 0

	r.offset++
	// fmt.Println("offset", r.offset)
	return bit, nil
}

func (r *BitReader) ReadNBit(n int) (int, error) {
	if n > 32 {
		return 0, errors.New("can't read more than 4byte")
	}
	out := 0
	code := 0
	for i := 0; i < n; i++ {
		bit, err := r.ReadBit()
		if err != nil {
			return 0, err
		}
		code = code << 1
		if bit {
			code = code | 0b1
		}
		if bit {
			out = out | (1 << i)
		}
	}
	// fmt.Println(strconv.FormatInt(int64(out), 2), strconv.FormatInt(int64(code), 2))
	return out, nil
}

func (r *BitReader) EndByte() {
	r.offset = 8
}

// Reads the next byte. Skips to the next byte-boundary
// https://datatracker.ietf.org/doc/html/rfc1951#page-11 3.2.4
func (r *BitReader) Read(p []byte) (n int, err error) {
	r.offset = 0
	return r.reader.Read(p)
}
