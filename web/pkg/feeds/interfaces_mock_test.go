// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

package feeds

import (
	gomock "github.com/golang/mock/gomock"
	podcast "github.com/mxpv/podcast"
	api "github.com/mxpv/podsync/web/pkg/api"
	reflect "reflect"
)

// Mockid is a mock of id interface
type Mockid struct {
	ctrl     *gomock.Controller
	recorder *MockidMockRecorder
}

// MockidMockRecorder is the mock recorder for Mockid
type MockidMockRecorder struct {
	mock *Mockid
}

// NewMockid creates a new mock instance
func NewMockid(ctrl *gomock.Controller) *Mockid {
	mock := &Mockid{ctrl: ctrl}
	mock.recorder = &MockidMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *Mockid) EXPECT() *MockidMockRecorder {
	return _m.recorder
}

// Generate mocks base method
func (_m *Mockid) Generate(feed *api.Feed) (string, error) {
	ret := _m.ctrl.Call(_m, "Generate", feed)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate indicates an expected call of Generate
func (_mr *MockidMockRecorder) Generate(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Generate", reflect.TypeOf((*Mockid)(nil).Generate), arg0)
}

// Mockstorage is a mock of storage interface
type Mockstorage struct {
	ctrl     *gomock.Controller
	recorder *MockstorageMockRecorder
}

// MockstorageMockRecorder is the mock recorder for Mockstorage
type MockstorageMockRecorder struct {
	mock *Mockstorage
}

// NewMockstorage creates a new mock instance
func NewMockstorage(ctrl *gomock.Controller) *Mockstorage {
	mock := &Mockstorage{ctrl: ctrl}
	mock.recorder = &MockstorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *Mockstorage) EXPECT() *MockstorageMockRecorder {
	return _m.recorder
}

// CreateFeed mocks base method
func (_m *Mockstorage) CreateFeed(feed *api.Feed) error {
	ret := _m.ctrl.Call(_m, "CreateFeed", feed)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateFeed indicates an expected call of CreateFeed
func (_mr *MockstorageMockRecorder) CreateFeed(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "CreateFeed", reflect.TypeOf((*Mockstorage)(nil).CreateFeed), arg0)
}

// GetFeed mocks base method
func (_m *Mockstorage) GetFeed(hashId string) (*api.Feed, error) {
	ret := _m.ctrl.Call(_m, "GetFeed", hashId)
	ret0, _ := ret[0].(*api.Feed)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeed indicates an expected call of GetFeed
func (_mr *MockstorageMockRecorder) GetFeed(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetFeed", reflect.TypeOf((*Mockstorage)(nil).GetFeed), arg0)
}

// Mockbuilder is a mock of builder interface
type Mockbuilder struct {
	ctrl     *gomock.Controller
	recorder *MockbuilderMockRecorder
}

// MockbuilderMockRecorder is the mock recorder for Mockbuilder
type MockbuilderMockRecorder struct {
	mock *Mockbuilder
}

// NewMockbuilder creates a new mock instance
func NewMockbuilder(ctrl *gomock.Controller) *Mockbuilder {
	mock := &Mockbuilder{ctrl: ctrl}
	mock.recorder = &MockbuilderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *Mockbuilder) EXPECT() *MockbuilderMockRecorder {
	return _m.recorder
}

// Build mocks base method
func (_m *Mockbuilder) Build(feed *api.Feed) (*podcast.Podcast, error) {
	ret := _m.ctrl.Call(_m, "Build", feed)
	ret0, _ := ret[0].(*podcast.Podcast)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Build indicates an expected call of Build
func (_mr *MockbuilderMockRecorder) Build(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Build", reflect.TypeOf((*Mockbuilder)(nil).Build), arg0)
}
