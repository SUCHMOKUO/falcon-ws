package util

import (
	"log"
	"os"
	"path/filepath"
)

var currentPath string

func init() {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalln("get current path error:", err)
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		log.Fatalln("get current path error:", err)
	}
	currentPath = filepath.Dir(exePath)
}

// GetCurrentPath return the path of
// the executable file.
func GetCurrentPath() string {
	return currentPath
}
