package main

import (
	"fmt"
	"flag"
	"os"
	"path/filepath"
)

const (
	FIND_FILES 		= 1 << 0
	FIND_DIRS 		= 1 << 1
	FIND_SYMLINKS 	= 1 << 2
)

type Finder interface {
	Print(in_path string)
}

type DirFinder struct {}
type FileFinder struct {}
type SymLinkFinder struct {}

func(dir DirFinder) Print(in_path string) {
	entries, err := os.ReadDir(in_path)
	if err != nil {
		return
	}
	for _, entry := range entries {
		full_path := filepath.Join(in_path, entry.Name())
		if entry.IsDir() {
			fmt.Println(full_path + "/")
			dir.Print(full_path)
		}
	}
}

func(file FileFinder) Print(in_path string) {
	entries, err := os.ReadDir(in_path)
	if err != nil {
		return
	}
	for _, entry := range entries {
		full_path := filepath.Join(in_path, entry.Name())
		if entry.IsDir() {
			file.Print(full_path)
		} else if entry.Type().IsRegular() {
			fmt.Println(full_path)
		}
	}
}

func(sl SymLinkFinder) Print(in_path string) {
	entries, err := os.ReadDir(in_path)
	if err != nil {
		return
	}
	for _, entry := range entries {
		full_path := filepath.Join(in_path, entry.Name())
		if entry.IsDir() {
			sl.Print(full_path)
		} else if entry.Type() & os.ModeSymlink != 0 {
			fmt.Println(full_path)
		}
	}
}

func print_data(in_path string) {
	entries, err := os.ReadDir(in_path)
	if err != nil {
		return
	}
	for _, entry := range entries {
		full_path := filepath.Join(in_path, entry.Name())
		if entry.IsDir() {
			fmt.Println(full_path + "/")
			print_data(full_path)
		} else {
			fmt.Println(full_path)
		}
	}
}

func main() {

	var mask int
	_ = mask

	f_file := flag.Bool("f", false, "file name to search")
	f_dir := flag.Bool("d", false, "directory name to search")
	f_symlink := flag.Bool("sl", false, "symlink name to search")
	// f_ext := flag.String("ext", "", "extantion files to search. works only if -f is specified")
	flag.Parse()

	var finder Finder
	
	args := flag.Args()
	if len(args) < 0 {
		fmt.Println("err: empty path to directory")
		os.Exit(1)
	}

	if *f_file {
		mask = mast & FIND_FILES
		finder = FileFinder{}
		finder.Print(args[0])
	}
	if *f_dir {
		mask = mast & FIND_DIRS
		finder = DirFinder{}
		finder.Print(args[0])
	}
	if *f_symlink {
		mask = mast & FIND_SYMLINKS
		finder = SymLinkFinder{}
		finder.Print(args[0])
	}



}
