package utils

import (
	"os"
	"os/exec"
	"path/filepath"
)

// GetCurrentPath get current path that running the execute file
func GetCurrentPath() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

// GetCurrentExecDir get directory of execute file
func GetCurrentExecDir() (dir string, err error) {
	var path string
	path, err = exec.LookPath(os.Args[0])
	if err != nil {
		// fmt.Printf("exec.LookPath(%s), err: %s\n", os.Args[0], err)
		return "", err
	}
	path, err = filepath.Abs(path)
	if err != nil {
		// fmt.Printf("filepath.Abs(%s), err: %s\n", path, err)
		return "", err
	}
	dir = filepath.Dir(path)
	return dir, nil
}
