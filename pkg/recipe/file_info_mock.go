package recipe

import (
	fs "io/fs"
	"time"
)

type mockFileInfo struct{}

func (*mockFileInfo) Name() string {
	return "DURIAN"
} // base name of the file
func (*mockFileInfo) Size() int64 {
	return 1
} // length in bytes for regular files; system-dependent for others
func (*mockFileInfo) Mode() fs.FileMode {
	return 1
}
func (*mockFileInfo) ModTime() time.Time {
	return time.Now()
}
func (*mockFileInfo) IsDir() bool {
	return false
}
func (*mockFileInfo) Sys() interface{} {
	return "APPLE"
}
