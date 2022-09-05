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

func (*fileMockTest) SetupTest() {
}

func TestFileMockTest(t *testing.T) {
	suite.Run(t, new(fileMockTest))
}

func (r *fileMockTest) Test_mockFileInfo() {
	mf := mockFileInfo{}

	r.Equal("DURIAN", mf.Name())
	r.Equal(int64(mockIntReturn), mf.Size())
	r.Equal(fs.FileMode(mockIntReturn), mf.Mode())
	r.Equal(time.Date(mockDateYear, mockDateMonth, mockDateDay, mockDateHour, mockDateMin, mockDateSec, mockDateNsec, time.UTC), mf.ModTime())
	r.Equal(false, mf.IsDir())
	r.Equal("APPLE", mf.Sys())
}
