package http

import (
	"errors"
	"strings"
)

type HTTPMethod int

const (
	GET = iota + 1
	PUT
	PATCH
	POST
)

func MethodFrom(str string) (HTTPMethod, error) {

	switch strings.ToUpper(str) {
	case "GET":
		return GET, nil
	case "PUT":
		return PUT, nil
	case "PATCH":
		return PATCH, nil
	case "POST":
		return POST, nil
	default:
		return 0, errors.New("unknown HttpMethod")
	}
}

func (method HTTPMethod) IsValid() bool {
	return method >= GET && method <= POST
}
