package utils

import (
	"os"
	"strings"
	"os/exec"
	"path/filepath"
)

func PathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

var appDir string
func GetAppDir() string {
	if appDir != "" {
		return  appDir
	}
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	appDir = filepath.Dir(path)

	return  appDir
}

func IsRelactivePath(path string)  bool {
	if strings.Index(path, ".") == 0 {
		return true
	}
	return false
}

func AbsPath(path string) string {
	if IsRelactivePath(path) {
		path = GetAppDir() + string(os.PathSeparator) + path
		return  path
	}
	return  path
}

