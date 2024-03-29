package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	Reset  = "\033[0m"
	Green = "\033[32m"
	Blue  = "\033[34m"
)

type Path struct {
	absPath, md5Hash string
	size             int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: /.duplicateFileHandlerGo <directory path>")
		return
	}

	root := os.Args[1]
	paths := getPaths(root, getFormat())

	sortingOrder := getSortingOrder()
	paths = sortPaths(paths, sortingOrder)
	printPaths(paths)

	duplicates := sortPaths(getDuplicates(paths), sortingOrder)
	wantsToCheckForDuplicates := wantsToCheckForDuplicates()
	if wantsToCheckForDuplicates && len(duplicates) == 0 {
		fmt.Println("No duplicates found. Bye!")
		return
	} else {
		printDuplicatePaths(duplicates)
	}

	if wantsToDeleteDuplicates() {
		deleteDuplicates(duplicates, getFileNumbersToDelete(duplicates))
	}

	return
}

func deleteDuplicates(duplicates []Path, numsToDelete []int) {
	var freedSpace int64
	for _, num := range numsToDelete {
		path := duplicates[num-1]
		freedSpace += path.size
		os.Remove(path.absPath)
	}

	fmt.Print(Blue + "\nTotal freed up space: ", freedSpace, " bytes\n\n" + Reset)
	return
}

func getFileNumbersToDelete(duplicates []Path) (fileNums []int) {
	fmt.Println(Blue + "\nEnter file numbers to delete (space separated integers):" + Reset)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		stringNums := strings.Fields(scanner.Text())
		if areFileNumsValid(duplicates, stringNums) {
			return convertToIntSlice(stringNums)
		}

		fmt.Print("\nPlease enter space separated integers only, corresponding to the duplicate files above.\n\nEnter file numbers to delete:")
	}
}

func areFileNumsValid(duplicates []Path, fileNums []string) bool {
	if len(fileNums) == 0 {
		return false
	}

	for _, num := range fileNums {
		fileNum, ok := strconv.Atoi(num)
		if ok != nil || (fileNum - 1 >= len(duplicates)) {
			return false
		}
	}

	return true
}

func convertToIntSlice(nums []string) (intNums []int) {
	for _, num := range nums {
		intNum, _ := strconv.Atoi(num)
		intNums = append(intNums, intNum)
	}

	return intNums
}

func wantsToDeleteDuplicates() bool {
	for {
	choice := readWord(Blue + "\nDelete files? (You will be prompted to select which files to delete)\n" + Reset)
		switch choice {
		case "yes":
			return true
		case "no":
			return false
		default:
			choice = readWord("\nEnter 'yes' or 'no'\n\n")
		}
	}
}

func printDuplicatePaths(paths []Path) {
	currentSize := int64(-1)
	currentHash := ""

	for i, path := range paths {
		if path.size != currentSize {
			fmt.Print(Green + "\n", path.size, " bytes\n" + Reset)
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

func getDuplicates(paths []Path) (duplicates []Path) {
	dupeMap := makeDupeMap(paths)
	for _, val := range dupeMap {
		if len(val) <= 1 {
			continue
		}
		duplicates = addPathsToDuplicate(duplicates, val)
	}

	return duplicates
}

func addPathsToDuplicate(duplicates, array []Path) (updatedDuplicates []Path) {
	for _, path := range array {
		duplicates = append(duplicates, path)
	}

	return duplicates
}

func makeDupeMap(paths []Path) (dupeMap map[string][]Path) {
	dupeMap = make(map[string][]Path)
	for _, path := range paths {
		_, ok := dupeMap[path.md5Hash]
		if !ok {
			pathArray := []Path{path}
			dupeMap[path.md5Hash] = pathArray
		} else {
			dupeMap[path.md5Hash] = append(dupeMap[path.md5Hash], path)
		}
	}

	return dupeMap
}

func printPaths(paths []Path) {
	currentSize := int64(-1)

	for _, path := range paths {
		if path.size != currentSize {
			fmt.Print(Green + "\n", path.size, Reset, " bytes\n")
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
	choice := readWord(Blue + "\nCheck for duplicates?\n" + Reset)
	for {
		switch choice {
		case "yes":
			return true
		case "no":
			return false
		default:
			choice = readWord("\nEnter 'yes' or 'no'\n\nCheck for duplicates?\n")
		}
	}
}

func getSortingOrder() (sortingOrder string) {
	selection := readInt("\nSize sorting options:\n1.Descending\n2.Ascending\n\n" + Blue + "Enter a sorting option:\n" + Reset)
	for {
		switch selection {
		case 1:
			return "Descending"
		case 2:
			return "Ascending"
		default:
			selection = readInt("\nPlease enter an integer representing one of the options\n\nEnter a sorting option:\n")
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
	fmt.Print(Blue + "Enter file format (enter nothing to search all files):\n" + Reset)
	fmt.Scanln(&line)
	return line
}
