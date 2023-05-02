package utils_collections

func Map[TSource any, TResult any](s []TSource, mapper func(TSource) TResult) []TResult {
	r := make([]TResult, 0)
	for _, t := range s {
		r = append(r, mapper(t))
	}
	return r
}
