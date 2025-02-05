package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

func create_gz(in_path string, wg *sync.WaitGroup, target_path string) {
	defer wg.Done()

	file, err := os.OpenFile(in_path, os.O_RDONLY, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", in_path, err)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting file info %s: %v\n", in_path, err)
		return
	}

	mtime := info.ModTime().Unix()
	archive_path := fmt.Sprintf("%s_%d.tar.gz",
		strings.TrimSuffix(filepath.Base(in_path), filepath.Ext(in_path)), mtime)

	if target_path != "" {
		archive_path = path.Join(target_path, archive_path)
	}

	if err := write_tar(file, archive_path, info); err != nil {
		fmt.Fprintf(os.Stderr, "Error archiving file %s: %v\n", in_path, err)
	}
}

func write_tar(file *os.File, archive_path string, info os.FileInfo) error {
	out_file, err := os.Create(archive_path)
	if err != nil {
		return fmt.Errorf("error creating archive %s: %w", archive_path, err)
	}
	defer out_file.Close()

	gz_writer := gzip.NewWriter(out_file)
	defer gz_writer.Close()
	tar_writer := tar.NewWriter(gz_writer)
	defer tar_writer.Close()

	header := &tar.Header{
		Name: filepath.Base(file.Name()),
		Size: info.Size(),
		Mode: int64(info.Mode()),
	}

	if err := tar_writer.WriteHeader(header); err != nil {
		return fmt.Errorf("error writing tar header: %w", err)
	}

	buf := make([]byte, 4096)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("error reading file %s: %w", file.Name(), err)
		}
		if _, err := tar_writer.Write(buf[:n]); err != nil {
			return fmt.Errorf("error writing to tar: %w", err)
		}
	}

	fmt.Printf("Archived: %s -> %s\n", file.Name(), archive_path)
	return nil
}

func main() {
	f_path := flag.String("a", "", "Directory to save archive")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: myRotate [-a target_path] file1, file2 ...]")
		os.Exit(1)
	}
	var wg sync.WaitGroup
	for _, arg := range args {
		wg.Add(1)
		go create_gz(arg, &wg, *f_path)
	}
	wg.Wait()
}
