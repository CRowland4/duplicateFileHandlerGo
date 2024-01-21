package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Directory is not specified")
		return
	}

	root := os.Args[1]
	printPaths(getPaths(root))

	return
}

func printPaths(paths []string) {
	for _, path := range paths {
		fmt.Println(path)
	}
}

func getPaths(root string) (paths []string) {
	_ = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			absPath, _ := filepath.Abs(path)
			paths = append(paths, absPath)
		}

		return nil
	})

	return paths
}
