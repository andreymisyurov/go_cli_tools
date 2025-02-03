package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	FIND_FILES    = 1 << 0
	FIND_DIRS     = 1 << 1
	FIND_SYMLINKS = 1 << 2
	FIND_EXT      = 1 << 3
)

func Finder(in_path string, mask uint8, in_ext string) {
	in_path, _ = filepath.Abs(in_path)
	entries, err := os.ReadDir(in_path)
	if err != nil {
		return
	}
	for _, entry := range entries {
		full_path := filepath.Join(in_path, entry.Name())
		if entry.IsDir() {
			if mask&FIND_DIRS != 0 {
				fmt.Println(full_path + "/")
			}
			Finder(full_path, mask, in_ext)
		} else if entry.Type().IsRegular() && mask&FIND_FILES != 0 {
			if mask&FIND_EXT == 0 || strings.HasSuffix(full_path, in_ext) {
				fmt.Println(full_path)
			}
		} else if entry.Type()&os.ModeSymlink != 0 && mask&FIND_SYMLINKS != 0 {
			_, err := filepath.EvalSymlinks(full_path)
			if err != nil {
				fmt.Println(full_path, "-> [broken]")
			} else {
				fmt.Println(full_path, "->", full_path)
			}
		}
	}
}

func main() {

	var mask uint8

	f_file := flag.Bool("f", false, "file name to search")
	f_dir := flag.Bool("d", false, "directory name to search")
	f_symlink := flag.Bool("sl", false, "symlink name to search")
	f_ext := flag.String("ext", "", "extantion files to search. works only if -f is specified")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("err: empty path to directory")
		os.Exit(1)
	}

	if *f_file {
		mask |= FIND_FILES
		if *f_ext != "" {
			mask |= FIND_EXT
		}
	}
	if *f_dir {
		mask |= FIND_DIRS
	}
	if *f_symlink {
		mask |= FIND_SYMLINKS
	}
	if mask == 0 {
		mask = FIND_DIRS | FIND_FILES | FIND_SYMLINKS
	}

	for _, arg := range args {
		Finder(arg, mask, *f_ext)
	}

}
