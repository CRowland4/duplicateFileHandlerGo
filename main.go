package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Path struct {
	absPath string
	size    int64
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

	return
}

func printPaths(paths []Path) {
	fmt.Println()
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

func getSortingOrder() (sortingOrder string) {
	prompt := "\nSize sorting options:\n1.Descending\n2.Ascending\n\nEnter a sorting option:\n"
	selection := readInt(prompt)

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
		paths = append(paths, Path{absPath: absPath, size: info.Size()})
		return nil
	})

	return paths
}

func readInt(prompt string) (num int) {
	fmt.Print(prompt)
	fmt.Scanln(&num)
	return num
}

func getFormat() (line string) {
	fmt.Print("Enter file format:\n")
	fmt.Scanln(&line)
	return line
}
