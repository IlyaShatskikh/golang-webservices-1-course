package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func printDir(out io.Writer, path string, printFiles bool, offset string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Fprintln(out, err)
	}

	// filter files if needed
	if !printFiles {
		n := 0
		for _, x := range files {
			if x.IsDir() {
				files[n] = x
				n++
			}
		}
		files = files[:n]
	}

	for i, file := range files {
		last := i == len(files)-1

		writeFileName(out, offset, file, last)

		if file.IsDir() {
			filePath := path + string(os.PathSeparator) + file.Name()

			if last {
				printDir(out, filePath, printFiles, offset+"\t")
			} else {
				printDir(out, filePath, printFiles, offset+"│\t")
			}
		}
	}
	return nil
}

func writeFileName(out io.Writer, offset string, file os.FileInfo, last bool) {
	var symbol string
	if last {
		symbol = "└───"
	} else {
		symbol = "├───"
	}

	var size string
	if !file.IsDir() {
		if file.Size() == 0 {
			size = " (empty)"
		} else {
			size = fmt.Sprintf(" (%vb)", file.Size())
		}
	}

	filename := offset + symbol + file.Name() + size
	fmt.Fprintln(out, filename)
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return printDir(out, path, printFiles, "")
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
