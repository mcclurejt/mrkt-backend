package common

import (
	"fmt"
	"reflect"
)

const DefaultTimeout = 60

type RequestOptions interface {
	ToQueryString() string
}

type ResultError struct {
	Result interface{}
	Error  error
}

// CollectResults - Returns an error if one occurred, otherwise returns an array of results
func CollectResults(ch chan ResultError, l int, target interface{}) error {
	t := reflect.ValueOf(target).Elem()
	for i := 0; i < l; i++ {
		tmp := <-ch
		err := tmp.Error
		if err != nil {
			return err
		}
		// append the array contained in result to target
		result := tmp.Result
		items := reflect.ValueOf(result)
		if items.Kind() == reflect.Slice {
			for j := 0; j < items.Len(); j++ {
				item := items.Index(j)
				if item.Kind() == reflect.Ptr || item.Kind() == reflect.Struct {
					t.Set(reflect.Append(t, item))
				}
			}
		}
	}
	return nil
}

type RouteNotRecognizedError struct {
	Route string
}

func (r RouteNotRecognizedError) Error() string {
	return fmt.Sprintf("Error, route '%s' not recognized\n", r.Route)
}

type OptionParseError struct {
	DesiredType reflect.Type
}

func (o OptionParseError) Error() string {
	return fmt.Sprintf("Error, unable to parse options to the desired type of '%s'\n", o.DesiredType)
}
