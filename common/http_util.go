package common

import "strings"

const (
	GET = "GET"
	POST = "POST"
	PUT = "PUT"
	DELETE = "DELETE"
)

func (r *Request) HeaderByName(name string) (string, bool) {
	for _, header := range r.Headers {
		if strings.ToLower(header.Key) == strings.ToLower(name) {
			return header.Value, true
		}
	}
	return "", false
}