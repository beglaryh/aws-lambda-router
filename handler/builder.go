package handler

import (
	"github.com/beglaryh/aws-lambda-router/context"
	"github.com/beglaryh/aws-lambda-router/http"
)

type RegisterBuilder struct {
	register *Handler
}

func Builder() *RegisterBuilder {
	var register Handler
	return &RegisterBuilder{
		register: &register,
	}
}

func (b *RegisterBuilder) Handler(handler func(context.Context) (http.Response, error)) *RegisterBuilder {
	b.register.Handler = handler
	return b
}

func (b *RegisterBuilder) ErrorHandler(errorHandler func(err error) (http.Response, error)) *RegisterBuilder {
	b.register.ErrorHandler = errorHandler
	return b
}

func (b *RegisterBuilder) MandatoryQueryParameters(params []string) *RegisterBuilder {
	b.register.MandatoryQueryParameters = params
	return b
}

func (b *RegisterBuilder) Build() Handler {
	return *b.register
}
