package utils_collections

import "errors"

func Map[TSource any, TResult any](s []TSource, mapper func(TSource) TResult) []TResult {
	r := make([]TResult, 0)
	for _, t := range s {
		r = append(r, mapper(t))
	}
	return r
}

func Find[TSource any](s []TSource, finder func(TSource) bool) (TSource, error) {
	var res TSource
	for _, item := range s {
		if finder(item) {
			return item, nil
		}
	}
	return res, errors.New("Not found")
}
