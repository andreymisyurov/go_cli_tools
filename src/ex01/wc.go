package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"unicode/utf8"
)

const (
	COUNT_LINE = 1 << 0
	COUNT_WORD = 1 << 1
	COUNT_CHAR = 1 << 2
)

var total uint64

func analize(path string, mask uint8, wg *sync.WaitGroup) {
	defer wg.Done()
	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		os.Exit(1)
	}
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	if mask&COUNT_WORD != 0 {
		scanner.Split(bufio.ScanWords)
	}

	var count uint64
	for scanner.Scan() {
		if mask&COUNT_CHAR != 0 {
			count += uint64(utf8.RuneCountInString(scanner.Text()))
		}
		count += 1
	}

	atomic.AddUint64(&total, count)
	fmt.Println(count, "\t", path)
}

func main() {
	f_line := flag.Bool("l", false, "counting line in file")
	f_word := flag.Bool("w", false, "counting words in file")
	f_char := flag.Bool("m", false, "counting chars in file")
	flag.Parse()
	var mask uint8
	var count uint8
	if *f_line == true {
		mask |= COUNT_LINE
		count++
	}
	if *f_word == true {
		mask |= COUNT_WORD
		count++
	}
	if *f_char == true {
		mask |= COUNT_CHAR
		count++
	}

	if count > 1 {
		os.Exit(1)
	}

	if mask == 0 {
		mask |= COUNT_WORD
	}
	args := flag.Args()
	var wg sync.WaitGroup
	for _, arg := range args {
		wg.Add(1)
		go analize(arg, mask, &wg)
	}
	wg.Wait()
	if len(args) > 1 {
		fmt.Println(atomic.LoadUint64(&total), "\t", "total")
	}
}
