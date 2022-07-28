package fs

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

const (
	ExistFileName    = "test.txt"
	NotExistFileName = "not_exist.txt"
)

func TestFilterByExt(t *testing.T) {
	files := []File{
		{
			Link: "/path/test.txt",
			Name: "test",
			Ext:  "txt",
			Size: 10,
		},
		{
			Link: "/path/test1.txt",
			Name: "test1",
			Ext:  "txt",
			Size: 10,
		},
		{
			Link: "/path/test2.jpg",
			Name: "test2",
			Ext:  "jpg",
			Size: 10,
		},
		{
			Link: "/path/test3.tst",
			Name: "test3",
			Ext:  "tst",
			Size: 10,
		},
	}
	result := map[string][]string{
		"txt": {"test", "test1"},
		"jpg": {"test2"},
		"tst": {"test3"},
	}
	for _, ext := range []string{"txt", "jpg", "tst"} {
		filtered := FilterByExt(files, ext)
		if len(filtered) != len(result[ext]) {
			t.Errorf("Incorrerct count of filtered elements for %s extension", ext)
		}
		for idx, file := range filtered {
			if file.Name != result[ext][idx] {
				t.Errorf("Incorrerct name of filtered element %v %s, want %s", ext, file.Name, result[ext][idx])
			}
		}
	}
}

func TestIsNotExist(t *testing.T) {
	content := []byte("test\n")
	if err := os.WriteFile(ExistFileName, content, 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(ExistFileName)

	for idx, file := range []string{ExistFileName, NotExistFileName} {
		want := []bool{false, true}
		res := IsNotExist(file)
		if res != want[idx] {
			t.Errorf("IsNotExist(%s) = %v; want %v", file, res, want[idx])
		}
	}

}

func TestIncrementFileName(t *testing.T) {
	filenames := []string{
		"test",
		"test.txt",
		"test.some.words.txt",
		"/home/user/test.txt",
		"c:\\Users\\Documents\\test.txt",
		"d:\\Пробел в пути\\Тест.txt",
	}
	results := []string{
		"test_1",
		"test_2.txt",
		"test.some.words_3.txt",
		"/home/user/test_4.txt",
		"c:\\Users\\Documents\\test_5.txt",
		"d:\\Пробел в пути\\Тест_6.txt",
	}

	for idx, filename := range filenames {
		res := IncrementFileName(filename, idx+1)
		if res != results[idx] {
			t.Errorf("IncrementFileName(%s, %d) = %v; want %v", filename, idx, res, results[idx])
		}
	}
}

func TestListDir(t *testing.T) {
	uploadDir := "upload/"
	filesQty := 5
	if err := os.Mkdir(uploadDir, 0777); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(uploadDir)
	for i := 0; i < filesQty; i++ {
		testStr := fmt.Sprintf("test %d\n", i)
		content := []byte(testStr)
		fileName := fmt.Sprintf("%stest_%d", uploadDir, i)
		if err := os.WriteFile(fileName, content, 0666); err != nil {
			t.Fatal(err)
		}
		defer os.Remove(fileName)
	}
	files, err := ListDir(uploadDir)
	if err != nil {
		t.Errorf("cannot list dir %s", uploadDir)
	}
	if len(files) != filesQty {
		t.Errorf("Incorrerct count of listed elements %d, want %d", len(files), filesQty)
	}
	for _, file := range files {
		if reflect.TypeOf(file) != reflect.TypeOf(File{}) {
			t.Errorf("Incorrerct type of listed elements %T, want %T", file, File{})
		}
	}
}

func TestSplitName(t *testing.T) {
	filenames := []string{
		"",
		"test",
		"test.txt",
		"test.some.words.txt",
	}
	results := [][]string{
		{"", ""},
		{"test", "txt"},
		{"test", "txt"},
		{"test.some.words", "txt"},
	}
	for idx, filename := range filenames {
		name, ext := SplitName(filename)
		if name != results[idx][0] && ext != results[idx][1] {
			t.Errorf("Incorrerct split, get '%s' name and '%s' ext, want '%s' name and '%s' ext", name, ext, results[idx][0], results[idx][1])
		}
	}
}
