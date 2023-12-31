// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	context "context"
	entity "url-shortener/pkg/app/urlshortener/entity"

	mock "github.com/stretchr/testify/mock"

	service "url-shortener/pkg/app/urlshortener/service"
)

// ShortenedURLService is an autogenerated mock type for the ShortenedURLService type
type ShortenedURLService struct {
	mock.Mock
}

type ShortenedURLService_Expecter struct {
	mock *mock.Mock
}

func (_m *ShortenedURLService) EXPECT() *ShortenedURLService_Expecter {
	return &ShortenedURLService_Expecter{mock: &_m.Mock}
}

// RetrieveShortenedURL provides a mock function with given fields: ctx, short
func (_m *ShortenedURLService) RetrieveShortenedURL(ctx context.Context, short string) (*entity.ShortenedURL, error) {
	ret := _m.Called(ctx, short)

	var r0 *entity.ShortenedURL
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*entity.ShortenedURL, error)); ok {
		return rf(ctx, short)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *entity.ShortenedURL); ok {
		r0 = rf(ctx, short)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.ShortenedURL)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, short)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ShortenedURLService_RetrieveShortenedURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RetrieveShortenedURL'
type ShortenedURLService_RetrieveShortenedURL_Call struct {
	*mock.Call
}

// RetrieveShortenedURL is a helper method to define mock.On call
//   - ctx context.Context
//   - short string
func (_e *ShortenedURLService_Expecter) RetrieveShortenedURL(ctx interface{}, short interface{}) *ShortenedURLService_RetrieveShortenedURL_Call {
	return &ShortenedURLService_RetrieveShortenedURL_Call{Call: _e.mock.On("RetrieveShortenedURL", ctx, short)}
}

func (_c *ShortenedURLService_RetrieveShortenedURL_Call) Run(run func(ctx context.Context, short string)) *ShortenedURLService_RetrieveShortenedURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *ShortenedURLService_RetrieveShortenedURL_Call) Return(shortenedURL *entity.ShortenedURL, err error) *ShortenedURLService_RetrieveShortenedURL_Call {
	_c.Call.Return(shortenedURL, err)
	return _c
}

func (_c *ShortenedURLService_RetrieveShortenedURL_Call) RunAndReturn(run func(context.Context, string) (*entity.ShortenedURL, error)) *ShortenedURLService_RetrieveShortenedURL_Call {
	_c.Call.Return(run)
	return _c
}

// ShortURL provides a mock function with given fields: ctx, url, opts
func (_m *ShortenedURLService) ShortURL(ctx context.Context, url string, opts *service.ShortURLOption) (*entity.ShortenedURL, error) {
	ret := _m.Called(ctx, url, opts)

	var r0 *entity.ShortenedURL
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *service.ShortURLOption) (*entity.ShortenedURL, error)); ok {
		return rf(ctx, url, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *service.ShortURLOption) *entity.ShortenedURL); ok {
		r0 = rf(ctx, url, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.ShortenedURL)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *service.ShortURLOption) error); ok {
		r1 = rf(ctx, url, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ShortenedURLService_ShortURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ShortURL'
type ShortenedURLService_ShortURL_Call struct {
	*mock.Call
}

// ShortURL is a helper method to define mock.On call
//   - ctx context.Context
//   - url string
//   - opts *service.ShortURLOption
func (_e *ShortenedURLService_Expecter) ShortURL(ctx interface{}, url interface{}, opts interface{}) *ShortenedURLService_ShortURL_Call {
	return &ShortenedURLService_ShortURL_Call{Call: _e.mock.On("ShortURL", ctx, url, opts)}
}

func (_c *ShortenedURLService_ShortURL_Call) Run(run func(ctx context.Context, url string, opts *service.ShortURLOption)) *ShortenedURLService_ShortURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(*service.ShortURLOption))
	})
	return _c
}

func (_c *ShortenedURLService_ShortURL_Call) Return(shortenedURL *entity.ShortenedURL, err error) *ShortenedURLService_ShortURL_Call {
	_c.Call.Return(shortenedURL, err)
	return _c
}

func (_c *ShortenedURLService_ShortURL_Call) RunAndReturn(run func(context.Context, string, *service.ShortURLOption) (*entity.ShortenedURL, error)) *ShortenedURLService_ShortURL_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewShortenedURLService interface {
	mock.TestingT
	Cleanup(func())
}

// NewShortenedURLService creates a new instance of ShortenedURLService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewShortenedURLService(t mockConstructorTestingTNewShortenedURLService) *ShortenedURLService {
	mock := &ShortenedURLService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
