package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Path struct {
	absPath, md5Hash string
	size             int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Directory is not specified")
		return
	}

	root := os.Args[1]
	paths := getPaths(root, getFormat())
	paths = sortPaths(paths, getSortingOrder())
	printPaths(paths)
	if wantsToCheckForDuplicates() {
		checkForDuplicates(paths)
	}

	return
}

func checkForDuplicates(paths []Path) {
	currentSize := int64(-1)
	currentHash := ""

	for i, path := range paths {
		if path.size != currentSize {
			fmt.Print("\n", path.size, " bytes\n")
			currentSize = path.size
		}
		if path.md5Hash != currentHash {
			fmt.Print("Hash: ", path.md5Hash, "\n")
			currentHash = path.md5Hash
		}
		fmt.Printf("%d. %s\n", i+1, path.absPath)
	}

	return
}

func printPaths(paths []Path) {
	currentSize := int64(-1)

	for _, path := range paths {
		if path.size != currentSize {
			fmt.Print("\n", path.size, " bytes\n")
			currentSize = path.size
		}
		fmt.Println(path.absPath)
	}

	return
}

func sortPaths(paths []Path, order string) (sortedPaths []Path) {
	sort.Slice(paths, func(i, j int) bool {
		if order == "Descending" {
			return paths[i].size > paths[j].size
		}
		return paths[i].size < paths[j].size
	})

	return paths
}

func wantsToCheckForDuplicates() (wantsCheck bool) {
	choice := readWord("\nCheck for duplicates?\n")
	for {
		switch choice {
		case "yes":
			return true
		case "no":
			return false
		default:
			choice = readWord("\nWrong Answer\n\nCheck for duplicates?\n")
		}
	}
}

func getSortingOrder() (sortingOrder string) {
	selection := readInt("\nSize sorting options:\n1.Descending\n2.Ascending\n\nEnter a sorting option:\n")
	for {
		switch selection {
		case 1:
			return "Descending"
		case 2:
			return "Ascending"
		default:
			selection = readInt("\nWrong option\n\nEnter a sorting option:\n")
		}
	}
}

func getPaths(root, format string) (paths []Path) {
	_ = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, format) {
			return nil
		}

		absPath, _ := filepath.Abs(path)
		paths = append(paths, Path{absPath: absPath, size: info.Size(), md5Hash: getMd5Hash(absPath)})
		return nil
	})

	return paths
}

func getMd5Hash(path string) (resultHash string) {
	file, _ := os.Open(path)
	defer file.Close()

	md5Hash := md5.New()
	io.Copy(md5Hash, file)
	resultHash = fmt.Sprintf("%x", md5Hash.Sum(nil))
	return resultHash
}

func readInt(prompt string) (num int) {
	fmt.Print(prompt)
	fmt.Scanln(&num)
	return num
}

func readWord(prompt string) (word string) {
	fmt.Print(prompt)
	fmt.Scan(&word)
	return word
}

func getFormat() (line string) {
	fmt.Print("Enter file format:\n")
	fmt.Scanln(&line)
	return line
}
