package recipe

import (
	"io/fs"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type fileMockTest struct {
	suite.Suite
}

func (t *fileMockTest) SetupTest() {
}

func TestFileMockTest(t *testing.T) {
	suite.Run(t, new(fileMockTest))
}

func (r *fileMockTest) Test_mockFileInfo() {
	mf := mockFileInfo{}

	r.Equal("DURIAN", mf.Name())
	r.Equal(int64(1), mf.Size())
	r.Equal(fs.FileMode(1), mf.Mode())
	r.Equal(time.Date(1998, time.April, 6, 0, 0, 0, 0, time.UTC), mf.ModTime())
	r.Equal(false, mf.IsDir())
	r.Equal("APPLE", mf.Sys())
}
