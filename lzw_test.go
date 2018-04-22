package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

const letters = "ABC"

var symbolTable map[string]uint16

func init() {
	//using a fixed-width implementation which can store 2^16 values
	symbolTable = make(map[string]uint16)
	//initialize dictionary with ascii characters corresponding to (A,B,C) (65,66,67)
	symbolTable[string(65)] = uint16(65)
	symbolTable[string(66)] = uint16(66)
	symbolTable[string(67)] = uint16(67)

	rand.Seed(time.Now().UnixNano())

}

//returns a random string of length n
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//creates a string of length n with random clumps of length maxClump
func RandStringClump(n int, maxClump int) string {
	b := make([]byte, n)
	for i := range b {
		if rand.Intn(100) == 5 && i > maxClump {
			letter := letters[rand.Intn(len(letters))]
			b[i] = letter
			for j := 0; j < maxClump; j++ {
				b[i-j] = letter
			}

		} else {
			b[i] = letters[rand.Intn(len(letters))]
		}
	}

	return string(b)
}

func BenchmarkRandString(b *testing.B) {
	totalLength := 10
	stringLength := make([]int, totalLength)
	for i := 0; i < len(stringLength); i++ {
		stringLength[i] = int(math.Pow(2, float64(i))) * 1000
	}

	for i := 0; i < len(stringLength); i++ {
		b.Run(strconv.Itoa(stringLength[i]), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				RandString(stringLength[i])
			}

		})
	}
}

func BenchmarkRandStringClump(b *testing.B) {
	totalLength := 8

	stringLength := make([]int, totalLength)

	for i := 0; i < len(stringLength); i++ {
		stringLength[i] = int(math.Pow(2, float64(i))) * 1000

	}

	for i := 0; i < len(stringLength); i++ {
		b.Run(strconv.Itoa(stringLength[i]), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				RandStringClump(stringLength[i], 20)
			}

		})
	}
}

func BenchmarkLZWCompress(b *testing.B) {

	count := 10
	fileSizeKB := make([]int, count)
	for i := 0; i < len(fileSizeKB); i++ {
		fileSizeKB[i] = int(math.Pow(2, float64(i))) * 1000
	}

	rfp := make([]*os.File, count)
	wfp := make([]*os.File, count)

	baseInputPath := []string{"testResults/in_", ""}
	baseOutputPath := []string{"testResults/out_", ""}

	for i := 0; i < count; i++ {
		size := strconv.Itoa(fileSizeKB[i])
		rfp[i], _ = os.OpenFile(strings.Join(baseInputPath, string(size)), os.O_RDWR|os.O_CREATE, 0755)
		rfp[i].WriteString(RandString(fileSizeKB[i]))
		rfp[i].Seek(0, 0)

	}

	for i := 0; i < count; i++ {
		size := strconv.Itoa(fileSizeKB[i])
		wfp[i], _ = os.OpenFile(strings.Join(baseOutputPath, size), os.O_RDWR|os.O_CREATE, 0755)
		b.Run(strconv.Itoa(fileSizeKB[i]), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				benchmarkLZWCompress(rfp[i], wfp[i], b)
			}

		})

		printCompressionRatio(rfp[i], wfp[i])
	}

}

func BenchmarkLZWCompressClump(b *testing.B) {
	count := 15
	fileSizeKB := make([]int, count)
	for i := 0; i < len(fileSizeKB); i++ {
		fileSizeKB[i] = int(math.Pow(2, float64(i))) * 1000
	}

	rfp := make([]*os.File, count)
	wfp := make([]*os.File, count)

	baseInputPath := []string{"testResults/in_clump_", ""}
	baseOutputPath := []string{"testResults/out_clump_", ""}

	for i := 0; i < count; i++ {
		size := strconv.Itoa(fileSizeKB[i])
		rfp[i], _ = os.OpenFile(strings.Join(baseInputPath, string(size)), os.O_RDWR|os.O_CREATE, 0755)
		rfp[i].WriteString(RandStringClump(fileSizeKB[i], 5))
		rfp[i].Seek(0, 0)

	}

	for i := 0; i < count; i++ {
		size := strconv.Itoa(fileSizeKB[i])
		wfp[i], _ = os.OpenFile(strings.Join(baseOutputPath, size), os.O_RDWR|os.O_CREATE, 0755)
		b.Run(strconv.Itoa(fileSizeKB[i]), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				benchmarkLZWCompress(rfp[i], wfp[i], b)
			}

		})

		printCompressionRatio(rfp[i], wfp[i])
	}

}

func benchmarkLZWCompress(r io.Reader, w io.Writer, b *testing.B) {

	app := ""
	result := make([]uint16, 0)
	dictSize := 3
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println(err)
	}
	for _, c := range bytes {
		wc := app + string(c)
		if _, ok := symbolTable[wc]; ok {
			app = wc
		} else {
			result = append(result, symbolTable[app])
			// Add wc to the dictionary.
			symbolTable[wc] = uint16(dictSize)
			dictSize++
			app = string(c)
		}
	}

	// Output the code for w.
	if app != "" {
		result = append(result, symbolTable[app])
	}

	error := binary.Write(w, binary.LittleEndian, result)
	if error != nil {
		fmt.Println(error)
	}

}

func printCompressionRatio(r *os.File, w *os.File) {
	ri, err := r.Stat()
	if err != nil {
		fmt.Println(err)
	}

	wi, err := w.Stat()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%f\n", float64(ri.Size())/float64(wi.Size()))

}
