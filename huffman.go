package main

import (
	"errors"
	"sort"
)

type huffmanDecoder struct {
	literal  *huffTree
	distance *huffTree
}

type huffLeaf struct {
	code int
	len  int
}

type huffTree struct {
	tree map[huffLeaf]int
}

func NewHuffmanDecoder(literalLen, distLen []int) *huffmanDecoder {
	literal := NewTree(literalLen)
	dist := NewTree(distLen)
	return &huffmanDecoder{
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

func createFixedHuffman() *huffmanDecoder {
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

func NewTree(lengths []int) *huffTree {
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

			m[huffLeaf{
				code: code,
				len:  len,
			}] = i
		}
	}
	return &huffTree{m}
}

func (t *huffTree) decode(in *BitReader) (int, error) {
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
		val, exists := t.tree[leaf]

		if exists {
			return val, nil
		}
	}
	// fmt.Println("bin:", strconv.FormatInt(int64(code), 2))
	return 0, errors.New("could not find code within 32bit")
}

func (h *huffmanDecoder) decodeLiteral(in *BitReader) (int, error) {
	return h.literal.decode(in)
}

func (h *huffmanDecoder) decodeDistance(in *BitReader) (int, error) {
	return h.distance.decode(in)
}

var dynamicSymbolOrder = []int{
	16, 17, 18, 0, 8, 7,
	9, 6, 10, 5, 11, 4,
	12, 3, 13, 2, 14, 1, 15,
}

// Creates tree to decode huffman-codes
func codeDecoder(numCodes int, in *BitReader) (*huffTree, error) {
	m := make(map[huffLeaf]int)

	var lengths []int
	for i := 0; i < numCodes; i++ {
		len, err := in.ReadNBit(3)
		if err != nil {
			return nil, err
		}
		lengths = append(lengths, len)
	}

	sortedLengths := NewSortByKey(dynamicSymbolOrder, lengths)
	sort.Sort(sortedLengths)

	maxLen := maxIntSlice(sortedLengths.value)
	counts := getLengthCounts(sortedLengths.value, maxLen)
	nextCode := make([]int, len(counts))

	code := 0
	counts[0] = 0
	for i := 1; i <= maxLen; i++ {
		code = (code + counts[i-1]) << 1
		nextCode[i] = code
	}

	for i, len := range sortedLengths.value {
		if len != 0 {
			code = nextCode[len]
			nextCode[len]++

			m[huffLeaf{
				code: code,
				len:  len,
			}] = sortedLengths.key[i]
			// fmt.Println(strconv.FormatInt(int64(code), 2), len, sortedLengths.key[i])
		}
	}
	return &huffTree{m}, nil
}

func decodeLenCodes(numLen int, tree *huffTree, in *BitReader) ([]int, error) {
	var lengths []int
	for len(lengths) < numLen {
		code, err := tree.decode(in)
		if err != nil {
			return nil, err
		}
		switch {
		case code < 16:
			lengths = append(lengths, code)
		case code == 16:
			count, err := in.ReadNBit(2)
			if err != nil {
				return nil, err
			}
			for j := 0; j < count+3; j++ {
				if len(lengths) == 0 {
					lengths = append(lengths, 0)
				} else {
					lengths = append(lengths, lengths[len(lengths)-1])
				}
			}

		case code == 17:
			count, err := in.ReadNBit(3)
			if err != nil {
				return nil, err
			}

			for j := 0; j < count+3; j++ {
				lengths = append(lengths, 0)
			}

		case code == 18:
			count, err := in.ReadNBit(7)
			if err != nil {
				return nil, err
			}

			for j := 0; j < count+11; j++ {
				lengths = append(lengths, 0)
			}
		}
	}
	return lengths, nil
}

func createDynamicHuffman(in *BitReader) (*huffmanDecoder, error) {
	HLIT, err := in.ReadNBit(5)
	if err != nil {
		return nil, err
	}
	HDIST, err := in.ReadNBit(5)
	if err != nil {
		return nil, err
	}
	HCLEN, err := in.ReadNBit(4)
	if err != nil {
		return nil, err
	}

	tree, err := codeDecoder(HCLEN+4, in)
	if err != nil {
		return nil, err
	}

	lengths, err := decodeLenCodes(HLIT+257+HDIST+1, tree, in)
	if err != nil {
		return nil, err
	}

	litLen := lengths[:HLIT+257]
	distLen := lengths[HLIT+257:]

	return NewHuffmanDecoder(litLen, distLen), nil
}
