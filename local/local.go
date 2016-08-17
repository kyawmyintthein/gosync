package local

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func Sync(file string, event int) {

	fmt.Println("Sync ------------------")
	fmt.Println(file)
	if event == 0 {
		removeLocalFile(file)
	}

	if event == 1 {
		createLocalFile(file)
	}

	if event == 2 {
		// copy file
		createLocalFile(file)
	}

	if event == 3 {
		createLocalFile(file)
	}
}

func getSyncDir() string {
	return "/Users/kyawmyintthein/gocode/src/gosync/target"
}

func getTargetDir(file string) string {
	targetDir := getSyncDir()
	targetFilePath := strings.Replace(file, "temp", targetDir, -1)
	d, _ := path.Split(targetFilePath)
	return d
}

func getTargetPath(file string) string {
	targetDir := getSyncDir()
	targetPath := strings.Replace(file, "temp", targetDir, -1)
	return targetPath
}

// remove file from sync folder
func removeLocalFile(file string) error {
	var err error
	targetFile := getTargetPath(file)
	fmt.Println(targetFile)
	if _, err = os.Stat(targetFile); os.IsNotExist(err) {
		return err
	}
	err = os.Remove(targetFile)
	return err
}

func createLocalFile(file string) error {
	var err error
	targetFile := getTargetPath(file)
	fmt.Println(targetFile)
	if _, err = os.Stat(targetFile); !os.IsNotExist(err) {
		err = os.Remove(targetFile)
		if err != nil {
			return err
		}
	}
	os.MkdirAll(getTargetDir(file), 0777)
	d1 := []byte("hello\ngo\n")
	err = ioutil.WriteFile(targetFile, d1, 0644)
	return err
}
