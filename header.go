package main

import (
	"bufio"
	"errors"
	"log"
)

type GHeader struct {
	ID1      byte
	ID2      byte
	CM       byte
	FLG      byte
	MTIME    [4]byte
	XFL      byte
	OS       byte
	fextra   *FExtra
	fname    *FName
	fcomment *FComment
	crc16    *CRC16
}

type FExtra struct {
	XLEN byte
	data []byte
}
type FName struct {
	name string
}

type FComment struct {
	comment string
}

type CRC16 struct {
	crc byte
}

func NewGHeader(r *bufio.Reader) (*GHeader, error) {
	var time [4]byte
	data := make([]byte, 10)

	// headerBuff := make([]byte, 10)
	n, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	if n < 10 {
		return nil, errors.New("could not read enough data")
	}

	copy(time[:], data[4:8])

	header := &GHeader{
		ID1:   data[0],
		ID2:   data[1],
		CM:    data[2],
		FLG:   data[3],
		MTIME: time,
		XFL:   data[8],
		OS:    data[9],
	}

	if header.FExtra() {
		skipFExtra(r)
	}
	if header.FName() {
		name, err := parseName(r)
		if err != nil {
			return nil, err
		}
		header.fname = name
	}
	if header.FComment() {
		comment, err := parseComment(r)
		if err != nil {
			return nil, err
		}
		header.fcomment = comment
	}
	if header.FHCRC() {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		header.crc16 = &CRC16{b}
	}

	return header, nil
}

func (h GHeader) FHCRC() bool {
	return h.FLG&0x2 != 0
}

func (h GHeader) FExtra() bool {
	return h.FLG&0x4 != 0
}

func (h GHeader) FName() bool {
	return h.FLG&0x8 != 0
}

func (h GHeader) FComment() bool {
	return h.FLG&0x10 != 0
}

func parseName(in *bufio.Reader) (*FName, error) {
	data, err := in.ReadBytes(0)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	s := string(data)
	return &FName{s}, nil
}

func parseComment(in *bufio.Reader) (*FComment, error) {
	data, err := in.ReadBytes(0)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	s := string(data)
	return &FComment{s}, nil
}

func skipFExtra(in *bufio.Reader) error {
	return errors.New("i'm wrong. xlen e to bytes")
	// len, err := in.ReadByte()
	// if err != nil {
	// 	return err
	// }
	// buff := make([]byte, len)
	// n, err := in.Read(buff)
	// if err != nil || n != int(len) {
	// 	return errors.New("fuck" + err.Error())
	// }
	// return nil
}
