package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

func OpenFile(config_file string) (fp *os.File) {
	return OpenFileWithPermission(config_file, 0777)
}
func OpenFileWithPermission(config_file string, permission os.FileMode) (fp *os.File) {

	//_home := os.Getenv("HOME")

	_home, _ := homedir.Dir()
	if strings.Contains(config_file, "~") {
		config_file = strings.Replace(config_file, "~", _home, -1)
	}
	fp, err := os.OpenFile(config_file, os.O_RDWR|os.O_CREATE, permission)
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}
	return fp
}

func RunTasks(command []string) {
	fmt.Println("-----called-----")
	for _, v := range command {
		fmt.Println(v)
	}
	fmt.Println("-----called-----")
}
