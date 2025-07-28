package exec

import "os"

var currentDir string

func init() {
	if dir, err := os.Getwd(); err == nil {
		currentDir = dir
	}
}

func GetCurrentDir() string {
	return currentDir
}
