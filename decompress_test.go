package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"testing"
)

// https://play.golang.org/p/S2s3e0nj4s :)
func compareReaders(result io.Reader, truth io.Reader) (bool, error) {
	buf1 := bufio.NewReader(result)
	buf2 := bufio.NewReader(truth)
	for {

		b1, err1 := buf1.ReadByte()
		b2, err2 := buf2.ReadByte()
		if err1 != nil && err1 != io.EOF {
			return false, err1
		}
		if err2 != nil && err2 != io.EOF {
			return false, err2
		}
		if err1 == io.EOF || err2 == io.EOF {
			return err1 == err2, nil
		}
		if b1 != b2 {
			return false, nil
		}
	}

}

func testDecompress(zipPath, truthPath string, t *testing.T) {
	truth, err := os.Open(truthPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer truth.Close()

	zip, err := os.Open(zipPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer truth.Close()

	var result bytes.Buffer
	err = decompress(zip, &result)
	if err != nil {
		t.Fatal(zipPath, ": ", err)
	}
	equal, err := compareReaders(&result, truth)
	if err != nil {
		t.Fatal(err)
	}
	if !equal {
		t.Fatal(zipPath, "Not equal to", truthPath)
	}
}

func TestDecompress(t *testing.T) {
	// testDecompress("testdata/short.txt.gz", "testdata/short.txt", t)
	// testDecompress("testdata/repeat.txt.gz", "testdata/repeat.txt", t)
	testDecompress("testdata/loremIpsum.txt.gz", "testdata/loremIpsum.txt", t)
}
