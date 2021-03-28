package main

import (
	"errors"
	"log"
)

type huffmanDecoder struct {
	literal  map[huffLeaf]int
	distance map[huffLeaf]int
}

type huffLeaf struct {
	code int
	len  int
}

func NewHuffmanDecoder(literalLen, distLen []int) huffmanDecoder {
	literal := NewTree(literalLen)
	dist := NewTree(distLen)
	return huffmanDecoder{
		literal:  literal,
		distance: dist,
	}
}

/*

   Lit Value    Bits        Codes
   ---------    ----        -----
     0 - 143     8          00110000 through
                            10111111
   144 - 255     9          110010000 through
                            111111111
   256 - 279     7          0000000 through
                            0010111
   280 - 287     8          11000000 through
                            11000111


	Distance 0-31 = 5bits
*/

func createFixedHuffman() huffmanDecoder {
	// Absolutely horrible, but it works.
	literalLen := make([]int, 288)

	for i := 0; i < 144; i++ {
		literalLen[i] = 8
	}
	for i := 144; i < 256; i++ {
		literalLen[i] = 9
	}
	for i := 256; i < 280; i++ {
		literalLen[i] = 7
	}
	for i := 280; i < 288; i++ {
		literalLen[i] = 8
	}

	distLen := make([]int, 32)
	for i := 0; i < 32; i++ {
		distLen[i] = 5
	}
	return NewHuffmanDecoder(literalLen, distLen)
}

func getLengthCounts(lengths []int, maxLen int) []int {
	counts := make([]int, maxLen+1)
	for _, len := range lengths {
		counts[len]++
	}
	return counts
}

func NewTree(lengths []int) map[huffLeaf]int {
	m := make(map[huffLeaf]int)

	maxLen := maxIntSlice(lengths)
	counts := getLengthCounts(lengths, maxLen)
	nextCode := make([]int, len(counts))

	code := 0
	counts[0] = 0
	for i := 1; i <= maxLen; i++ {
		code = (code + counts[i-1]) << 1
		nextCode[i] = code
	}

	for i, len := range lengths {
		if len != 0 {
			code = nextCode[len]
			nextCode[len]++
			_, exists := m[huffLeaf{
				code: code,
				len:  len,
			}]
			if exists {
				log.Fatalln("Overwriting huffmancode")
			}
			m[huffLeaf{
				code: code,
				len:  len,
			}] = i
		}
	}
	return m
}

func (h *huffmanDecoder) decodeLiteral(in *BitReader) (int, error) {

	code := 0
	for i := 0; i < 10; i++ {
		bit, err := in.ReadBit()
		// fmt.Println(bit)
		if err != nil {
			return 0, err
		}
		code = code << 1
		if bit {
			code = code | 0b1
		}
		leaf := huffLeaf{
			code: code,
			len:  i + 1,
		}
		val, exists := h.literal[leaf]

		if exists {
			// fmt.Println(strconv.FormatInt(int64(code), 2))
			return val, nil
		}
	}
	return 0, errors.New("could not find code within 32bit")
}

func (h *huffmanDecoder) decodeDistance(in *BitReader) (int, error) {
	code := 0
	for i := 0; i < 32; i++ {
		bit, err := in.ReadBit()
		if err != nil {
			return 0, err
		}
		code = code << 1
		if bit {
			code = code | 0b1
		}

		leaf := huffLeaf{
			code: code,
			len:  i + 1,
		}
		val, exists := h.distance[leaf]
		if exists {
			return val, nil
		}
	}

	return 0, errors.New("could not find code within 32bit")
}
