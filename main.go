package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type BufPos int
type Literal []byte

const EndOfBlock = 256

func writeBlock(in *BitReader, out io.Writer, decoder huffmanDecoder) error {
	ring := NewRingBuffer(3 * 1024)
	buf := make([]byte, 2)

	writer := io.MultiWriter(out, ring)

	for {
		val, err := decoder.decodeLiteral(in)
		// fmt.Println("val", val, err)
		if err != nil {
			return err
		}
		if val < 256 {

			binary.BigEndian.PutUint16(buf, uint16(val))
			writer.Write(buf[1:])
		} else if val == 256 {
			// fmt.Println("\n\nVal == 256, breaking")
			break
		} else {
			// fmt.Println("SMART")
			lenCode := lengthCodes[val]

			addedLen, err := in.ReadNBit(lenCode.bitCount)
			if err != nil {
				return err
			}
			distVal, err := decoder.decodeDistance(in)
			if err != nil {
				return err
			}
			distCode := distanceCodes[distVal]
			addedDist, err := in.ReadNBit(distCode.bitCount)
			if err != nil {
				return err
			}

			for i := 0; i < addedLen+lenCode.len; i++ {
				data, err := ring.LookBack(addedDist+distCode.len, 1)
				if err != nil {
					return err
				}
				writer.Write(data)
			}

		}
	}
	return nil
}

func writeDynamicBlock(in *BitReader, out io.Writer) error {
	decoder := createDynamicHuffman(in)
	return writeBlock(in, out, decoder)

}

func writeFixedBlock(in *BitReader, out io.Writer) error {
	decoder := createFixedHuffman()
	return writeBlock(in, out, decoder)
}

func writeUncompressedBlock(in *BitReader, out io.Writer) error {
	var lenBuf [2]byte
	var nlenBuf [2]byte

	_, err := io.ReadFull(in, lenBuf[:])
	if err != nil {
		return err
	}

	_, err = io.ReadFull(in, nlenBuf[:])
	if err != nil {
		return err
	}

	length := binary.LittleEndian.Uint16(lenBuf[:])
	_, err = io.CopyN(out, in, int64(length))
	if err != nil {
		return err
	}

	return nil
}

func decompressBlock(in *BitReader, out io.Writer) (bool, error) {
	var err error

	header, err := in.ReadNBit(3)
	if err != nil {
		return false, err
	}

	bHeader := NewBlockHeader(byte(header))
	fmt.Println("Header", bHeader.Bfinal, bHeader.BType)
	switch bHeader.BType {
	case blockNoCompression:
		err = writeUncompressedBlock(in, out)
	case blockFixed:
		err = writeFixedBlock(in, out)
	case blockDynamic:
		err = writeDynamicBlock(in, out)
	default:
		return bHeader.Bfinal, errors.New("cant decode block-type:" + fmt.Sprint(bHeader.BType))
	}
	if err != nil {
		return bHeader.Bfinal, err
	}

	return bHeader.Bfinal, nil
}

func decompress(in io.Reader, out io.Writer) error {

	r := bufio.NewReader(in)

	_, err := NewGHeader(r)
	if err != nil {
		return err
	}
	bitReader := NewBitReader(r)
	final := false
	for !final {
		final, err = decompressBlock(bitReader, out)
		if err != nil {
			return err
		}
	}

	return nil
}

func compress(in io.Reader, out io.Writer) error {
	return errors.New("compress not implemented")
}

func main() {
	var input io.Reader
	var err error

	decode := flag.Bool("d", false, "Decompress file")
	flag.Parse()
	files := flag.Args()

	if len(files) != 1 {
		input = os.Stdin
	} else {

		file, err := os.Open(files[0])
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()
		input = file
	}

	if *decode {
		err = decompress(input, os.Stdout)
	} else {
		err = compress(input, os.Stdout)
	}
	if err != nil {
		log.Fatalln(err)
	}
}
