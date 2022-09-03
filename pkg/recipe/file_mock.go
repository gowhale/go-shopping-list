package recipe

import (
	fs "io/fs"
	"time"
)

const (
	mockIntReturn = 1

	mockDateYear  = 1998
	mockDateMonth = time.April
	mockDateDay   = 26
	mockDateHour  = 0
	mockDateMin   = 0
	mockDateSec   = 0
	mockDateNsec  = 0
)

type mockFileInfo struct{}

func (*mockFileInfo) Name() string {
	return "DURIAN"
}
func (*mockFileInfo) Size() int64 {
	return mockIntReturn
} // length in bytes for regular files; system-dependent for others
func (*mockFileInfo) Mode() fs.FileMode {
	return mockIntReturn
}
func (*mockFileInfo) ModTime() time.Time {
	return time.Date(mockDateYear, mockDateMonth, mockDateDay, mockDateHour, mockDateMin, mockDateSec, mockDateNsec, time.UTC)
}
func (*mockFileInfo) IsDir() bool {
	return false
}
func (*mockFileInfo) Sys() interface{} {
	return "APPLE"
}
