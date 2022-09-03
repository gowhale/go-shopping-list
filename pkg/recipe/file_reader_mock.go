// Code generated by mockery v1.0.0. DO NOT EDIT.

package recipe

import (
	fs "io/fs"

	mock "github.com/stretchr/testify/mock"
)

// MockFileReader is an autogenerated mock type for the FileReader type
type MockFileReader struct {
	mock.Mock
}

// GetPopularity provides a mock function with given fields: recipeName
func (_m *MockFileReader) GetPopularity(recipeName string) (int, error) {
	ret := _m.Called(recipeName)

	var r0 int
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(recipeName)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(recipeName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IncrementPopularity provides a mock function with given fields: recipeName
func (_m *MockFileReader) IncrementPopularity(recipeName string) error {
	ret := _m.Called(recipeName)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(recipeName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LoadPopularityFile provides a mock function with given fields:
func (_m *MockFileReader) LoadPopularityFile() (Popularity, error) {
	ret := _m.Called()

	var r0 Popularity
	if rf, ok := ret.Get(0).(func() Popularity); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(Popularity)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadRecipeFile provides a mock function with given fields: fileName
func (_m *MockFileReader) LoadRecipeFile(fileName fs.FileInfo) (Recipe, error) {
	ret := _m.Called(fileName)

	var r0 Recipe
	if rf, ok := ret.Get(0).(func(fs.FileInfo) Recipe); ok {
		r0 = rf(fileName)
	} else {
		r0 = ret.Get(0).(Recipe)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(fs.FileInfo) error); ok {
		r1 = rf(fileName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MarshallJSON provides a mock function with given fields: pop
func (_m *MockFileReader) MarshallJSON(pop Popularity) ([]byte, error) {
	ret := _m.Called(pop)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(Popularity) []byte); ok {
		r0 = rf(pop)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(Popularity) error); ok {
		r1 = rf(pop)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadFile provides a mock function with given fields: filePath
func (_m *MockFileReader) ReadFile(filePath string) ([]byte, error) {
	ret := _m.Called(filePath)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(filePath)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(filePath)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadRecipeDirectory provides a mock function with given fields:
func (_m *MockFileReader) ReadRecipeDirectory() ([]fs.FileInfo, error) {
	ret := _m.Called()

	var r0 []fs.FileInfo
	if rf, ok := ret.Get(0).(func() []fs.FileInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]fs.FileInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnmarshallJSONToPopularity provides a mock function with given fields: file
func (_m *MockFileReader) UnmarshallJSONToPopularity(file []byte) (Popularity, error) {
	ret := _m.Called(file)

	var r0 Popularity
	if rf, ok := ret.Get(0).(func([]byte) Popularity); ok {
		r0 = rf(file)
	} else {
		r0 = ret.Get(0).(Popularity)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(file)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnmarshallJSONToRecipe provides a mock function with given fields: file
func (_m *MockFileReader) UnmarshallJSONToRecipe(file []byte) (Recipe, error) {
	ret := _m.Called(file)

	var r0 Recipe
	if rf, ok := ret.Get(0).(func([]byte) Recipe); ok {
		r0 = rf(file)
	} else {
		r0 = ret.Get(0).(Recipe)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(file)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteFile provides a mock function with given fields: newFile
func (_m *MockFileReader) WriteFile(newFile []byte) error {
	ret := _m.Called(newFile)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(newFile)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WritePopularityFile provides a mock function with given fields: pop
func (_m *MockFileReader) WritePopularityFile(pop Popularity) error {
	ret := _m.Called(pop)

	var r0 error
	if rf, ok := ret.Get(0).(func(Popularity) error); ok {
		r0 = rf(pop)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
