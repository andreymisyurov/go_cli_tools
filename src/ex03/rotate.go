package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func create_gz(in_path string, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.OpenFile(in_path, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer file.Close()

	info, _ := file.Stat()
	mtime := info.ModTime().Unix()
	archive_path := fmt.Sprintf("%s_%d.tar.gz", filepath.Base(in_path), mtime)
	out_file, err := os.Create(archive_path)
	if err != nil {
		return
	}
	defer out_file.Close()

	gz_writter := gzip.NewWriter(out_file)
	defer gz_writter.Close()
	tar_writer := tar.NewWriter(gz_writter)
	defer tar_writer.Close()
	header := &tar.Header{
		Name: filepath.Base(in_path),
		Size: info.Size(),
		Mode: 0600,
	}

	if err := tar_writer.WriteHeader(header); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing header: %v\n", err)
		return
	}

	buf := make([]byte, 4096)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if _, err := tar_writer.Write(buf[:n]); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing tar: %v\n", err)
			return
		}
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		os.Exit(1)
	}
	var wg sync.WaitGroup
	for _, arg := range args {
		wg.Add(1)
		go create_gz(arg, &wg)
	}
	wg.Wait()
	fmt.Println(args)
}
