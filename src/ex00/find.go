package main

import (
	"fmt"
	"flag"
	"os"
	"path/filepath"
)

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
	f_file := flag.Bool("f", false, "file name to search")
	f_dir := flag.Bool("d", false, "directory name to search")
	f_symlink := flag.Bool("sl", false, "symlink name to search")
	// f_ext := flag.String("ext", "", "extantion files to search. works only if -f is specified")
	flag.Parse()

	count := 0
	if *f_file {
		count += 1
	}
	if *f_dir {
		count += 1
	}
	if *f_symlink {
		count += 1
	}

	if count > 1 {
		fmt.Println("err: you should call only one flag")
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) < 0 {
		fmt.Println("err: empty path to directory")
		os.Exit(1)
	}

	print_data(args[0])

	// fmt.Println(*f_file)
	// fmt.Println(*f_dir)
	// fmt.Println(*f_symlink)
	// fmt.Println(path)

}
