package helper

import (
	"path/filepath"
	"runtime"
)

func GetProjectDirectory() string {
	_, filename, _, ok := runtime.Caller(1)

	if !ok {
		panic("Cannot resolve current file path - why ?!")
	}

	dir := filepath.Clean(filepath.Dir(filename) + "../../../")

	return dir
}

func ResolveProjectFile(file string) string {
	return filepath.Join(GetProjectDirectory(), file)
}
