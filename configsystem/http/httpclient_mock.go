package http

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type HttpClientMock struct{
	mock.Mock
}

func (httpClient *HttpClientMock) Get(url string) (resp *http.Response, err error){
	httpClient.Called(url)
	return nil,nil
}

func (httpClient *HttpClientMock) Head(url string) (resp *http.Response, err error){
	httpClient.Called(url)
	return nil,nil
}