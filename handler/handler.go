package handler

import (
	"github.com/beglaryh/aws-lambda-router/context"
	"github.com/beglaryh/aws-lambda-router/http"
)

type Handler struct {
	Handler                  func(context.Context) (http.Response, error)
	ErrorHandler             func(err error) (http.Response, error)
	MandatoryQueryParameters []string
}
