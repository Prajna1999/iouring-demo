package main

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

// panic
func assert(b bool) {
	if !b {
		panic("assert")
	}
}

const BUFFER_SIZE = 4096 //4KiB

func readNBytes(fn string, n int) []byte {
	f, err := os.Open(fn)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	data := make([]byte, 0, n)

	var buffer = make([]byte, BUFFER_SIZE)

	for len(data) < n {
		read, err := f.Read(buffer)
		if err != nil {
			panic(err)
		}
		data = append(data, buffer[:read]...)
	}

	assert(len(data) == n)
	return data

}

// benchmarking function
func benchmark(name string, data []byte, fn func(*os.File)) {
	fmt.Printf("%s", name)
	f, err := os.OpenFile("out.bin", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

	if err != nil {
		panic(err)
	}
	t1 := time.Now()

	fn(f)
	s := time.Since(t1).Seconds()
	fmt.Printf(",%f,%f\n", s, float64(len(data))/s)
	if err := f.Close(); err != nil {
		panic(err)
	}
	assert(bytes.Equal(readNBytes("out.bin", len(data)), data))
}

func main() {
	size := 104857600 //100MiB
	data := readNBytes("/dev/random", size)

	const RUNS = 20
	for i := 0; i < RUNS; i++ {
		benchmark("blocking", data, func(f *os.File) {
			for i := 0; i < len(data); i += BUFFER_SIZE {
				size := min(BUFFER_SIZE, len(data)-i)
				n, err := f.Write(data[i : i+size]) //number of bytes written
				if err != nil {
					panic(err)
				}
				assert(n == BUFFER_SIZE)
			}
		})
	}
}
