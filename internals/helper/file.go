package helper

import "path"

func GetProjectDirectory() string {
	return path.Clean("../../")
}

func ResolveProjectFile(file string) string {
	return path.Join(GetProjectDirectory(), file)
}
