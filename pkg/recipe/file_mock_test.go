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
	r.Equal(mf.Size(), int64(1))
	r.Equal(mf.Mode(), fs.FileMode(1))
	r.Equal(mf.ModTime(), time.Date(1998, time.April, 6, 0, 0, 0, 0, time.UTC))
	r.Equal(mf.IsDir(), false)
	r.Equal(mf.Sys(), "APPLE")
}
