package api

import "net/url"

func queryToParams(q url.Values) map[string]string {
	res := make(map[string]string)
	for k, v := range q {
		if len(v) > 0 {
			res[k] = v[0]
		}
	}
	return res
}
