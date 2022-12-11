package utils

import "os"

func FileExists(path string) bool {
	if file, err := os.Stat(path); os.IsNotExist(err) || file.IsDir() {
		return false
	}
	return true
}

func DirExists(path string) bool {
	if dir, err := os.Stat(path); os.IsNotExist(err) || !dir.IsDir() {
		return false
	}
	return true
}
