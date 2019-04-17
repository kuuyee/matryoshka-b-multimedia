package test

import (
	"path"
	"path/filepath"
	"runtime"
)

// GetProjectDir returns the currect apsolute path of this project
func GetProjectDir() string {
	_, f, _, _ := runtime.Caller(0)
	projectDir, _ := filepath.Abs(path.Join(filepath.Dir(f), "../"))
	return projectDir
}
