package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

const (
	COUNT_LINE = 1 << 0
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
	scanner := bufio.NewScanner(file)

	var count uint64
	for scanner.Scan() {
		if mask&COUNT_LINE != 0 {
			count += 1
		}
	}
	atomic.AddUint64(&total, count)
	fmt.Println(count, path)
}

func main() {
	f_line := flag.Bool("l", false, "counting line in file")
	flag.Parse()
	var mask uint8
	if *f_line == true {
		mask |= COUNT_LINE
	}
	args := flag.Args()
	var wg sync.WaitGroup
	for _, arg := range args {
		wg.Add(1)
		go analize(arg, mask, &wg)
	}
	wg.Wait()
	fmt.Println(atomic.LoadUint64(&total), "total")
}
