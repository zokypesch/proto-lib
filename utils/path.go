package utils

import (
	"os"
	"path/filepath"
	"strings"
)

//GetPath for get path
func GetPath() string {
	gopath := os.Getenv("GOPATH")

	gopath += "/src/"
	filePath, _ := filepath.Abs("./")
	filePath = strings.Replace(filePath, gopath, "", -1)

	return filePath
}

//GetFullPath for get path
func GetFullPath() string {
	filePath, _ := filepath.Abs("./")

	return filePath
}
