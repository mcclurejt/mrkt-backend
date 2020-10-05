package base

const DefaultTimeout = 60

type RequestOptions interface {
	ToQueryString() string
}
