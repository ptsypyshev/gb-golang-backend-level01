package fs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type File struct {
	Link string
	Name string
	Ext  string
	Size int64
}

func SplitName(filename string) (name, ext string) {
	splitName := strings.Split(filename, ".")
	switch len(splitName) {
	case 1:
		return filename, ""
	case 2:
		return splitName[0], splitName[1]
	default:
		lastElem := len(splitName) - 1
		return strings.Join(splitName[:lastElem], "."), splitName[lastElem]
	}
}

func ListDir(root string) (files []File, e error) {
	errWalk := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() {
			name, ext := SplitName(info.Name())
			baseFileName := filepath.Base(path)
			fileLink := "files/" + baseFileName
			files = append(files, File{
				Link: fileLink,
				Name: name,
				Ext:  ext,
				Size: info.Size(),
			})
		}
		return nil
	})

	if errWalk != nil {
		return nil, fmt.Errorf("error walking the path %q: %v", root, errWalk)
	}

	return files, nil
}

func FilterByExt(files []File, ext string) (res []File) {
	for _, file := range files {
		if file.Ext == ext {
			res = append(res, file)
		}
	}
	return
}

func IsNotExist(filePath string) bool {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return true
	}
	return false
}

func IncrementFileName(filePath string, addNum int) string {
	newTail := "_" + strconv.Itoa(addNum)
	nameSlice := strings.Split(filePath, ".")
	if len(nameSlice) > 1 {
		nameSlice[len(nameSlice)-2] += newTail
		return strings.Join(nameSlice, ".")
	}
	return filePath + newTail
}
