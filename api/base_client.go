package api

import "net/http"

type requestOptions interface {
	ToQueryString() string
}

type baseClient interface {
	call(options requestOptions) (*http.Response, error)
}
