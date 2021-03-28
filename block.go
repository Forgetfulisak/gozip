package main

type blockHeader struct {
	Bfinal bool
	BType  blockType
}

type blockType int

const (
	blockNoCompression blockType = iota
	blockFixed
	blockDynamic
)

var btypeCodes = map[int]blockType{
	0b00: blockNoCompression,
	0b01: blockFixed,
	0b10: blockDynamic,
}

func NewBlockHeader(data byte) blockHeader {
	btype := btypeCodes[int((data>>1)&0x3)]
	return blockHeader{
		Bfinal: data&0x1 != 0,
		BType:  btype,
	}
}
