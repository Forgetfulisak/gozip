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

func writeBlock(in *BitReader, out io.Writer, decoder *huffmanDecoder, lookBackBuffer *RingBuffer) error {
	buf := make([]byte, 2)

	writer := io.MultiWriter(out, lookBackBuffer)

	for {
		val, err := decoder.decodeLiteral(in)
		if err != nil {
			return err
		}

		if val < 256 {

			binary.BigEndian.PutUint16(buf, uint16(val))
			writer.Write(buf[1:])
			// fmt.Println("written", val, buf[1:], buf)
		} else if val == 256 {
			// fmt.Println("\n\nVal == 256, breaking")
			break
		} else {
			// fmt.Println("SMART")
			lenCode, ok := lengthCodes[val]
			if !ok {
				panic("not ok lencodes!")
			}

			addedLen, err := in.ReadNBit(lenCode.bitCount)
			if err != nil {
				return err
			}
			distVal, err := decoder.decodeDistance(in)
			if err != nil {
				return err
			}
			distCode, ok := distanceCodes[distVal]
			if !ok {
				panic("not ok distancecodes!")
			}

			addedDist, err := in.ReadNBit(distCode.bitCount)
			if err != nil {
				return err
			}

			for i := 0; i < addedLen+lenCode.len; i++ {

				data, err := lookBackBuffer.LookBack(addedDist+distCode.len, 1)
				if err != nil {
					return err
				}

				writer.Write(data)
			}

		}
	}
	return nil
}

func writeDynamicBlock(in *BitReader, out io.Writer, lookBackBuffer *RingBuffer) error {
	decoder, err := createDynamicHuffman(in)
	if err != nil {
		return err
	}
	return writeBlock(in, out, decoder, lookBackBuffer)
	// return errors.New("writeDynamicBlock not implemented")

}

func writeFixedBlock(in *BitReader, out io.Writer, lookBackBuffer *RingBuffer) error {
	decoder := createFixedHuffman()
	return writeBlock(in, out, decoder, lookBackBuffer)
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

func decompressBlock(in *BitReader, out io.Writer, lookBackBuffer *RingBuffer) (bool, error) {
	var err error

	header, err := in.ReadNBit(3)
	if err != nil {
		return false, err
	}

	bHeader := NewBlockHeader(byte(header))

	switch bHeader.BType {
	case blockNoCompression:
		// TODO: Test
		// Couldn't get gzip to create uncompressed block
		err = writeUncompressedBlock(in, out)
	case blockFixed:
		err = writeFixedBlock(in, out, lookBackBuffer)
	case blockDynamic:
		err = writeDynamicBlock(in, out, lookBackBuffer)
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

	lookBackBuffer := NewRingBuffer(32 * 1024)

	_, err := NewGHeader(r)
	if err != nil {
		return err
	}
	bitReader := NewBitReader(r)
	final := false
	for !final {
		final, err = decompressBlock(bitReader, out, lookBackBuffer)
		if err != nil {
			return err
		}
	}

	return nil
}

func compress(in io.Reader, out io.Writer) error {
	// outBuf := make([]byte, 10)
	// outBuf[0] = 0x1f
	// outBuf[1] = 0x0b
	// out.Write(outBuf)
	// r := bufio.NewReader(in)

	// for {

	// 	_, err := r.Peek(1)
	// 	if err == io.EOF {
	// 		break
	// 	} else if err != nil {
	// 		return err
	// 	}
	// 	header := []byte{0b00000000}

	// 	len := uint16(r.Buffered())
	// 	xlen := len ^ 0xffff

	// 	buf := make([]byte, 2)
	// 	binary.LittleEndian.PutUint16(buf, len)
	// 	out.Write(buf)
	// 	binary.LittleEndian.PutUint16(buf, xlen)
	// 	out.Write(buf)

	// 	io.CopyN(out, r, int64(len))
	// }
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
