package router

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/beglaryh/aws-lambda-router/context"
	"github.com/beglaryh/aws-lambda-router/handler"
	"github.com/beglaryh/aws-lambda-router/http"
)

type Router struct {
	gets    map[string]handler.Handler
	puts    map[string]handler.Handler
	patches map[string]handler.Handler
	posts   map[string]handler.Handler
}

func New() *Router {
	return &Router{
		gets: map[string]handler.Handler{},
	}
}

func (r *Router) RegisterGet(path string, handler handler.Handler) error {
	return r.register(http.GET, path, handler)
}

func (r *Router) RegisterPut(path string, handler handler.Handler) error {
	return r.register(http.PUT, path, handler)
}

func (r *Router) RegisterPatch(path string, handler handler.Handler) error {
	return r.register(http.PATCH, path, handler)
}

func (r *Router) RegisterPost(path string, handler handler.Handler) error {
	return r.register(http.POST, path, handler)
}

func (r *Router) Route(event events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	httpMethod, err := http.MethodFrom(event.HTTPMethod)

	if err != nil {
		return defaultErrorResponse()
	}

	switch httpMethod {
	case http.GET:
		return route(r.gets, event)
	case http.PUT:
		return route(r.puts, event)
	case http.PATCH:
		return route(r.patches, event)
	case http.POST:
		return route(r.posts, event)
	default:
		response := events.APIGatewayProxyResponse{
			StatusCode: 405,
		}
		return response
	}

}
func route(methodPaths map[string]handler.Handler, event events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	handler, ok := methodPaths[event.Resource]
	if !ok {
		response := events.APIGatewayProxyResponse{
			StatusCode: 404,
		}
		return response
	}
	validationResponse, pass := validateQueryParameters(handler.MandatoryQueryParameters, event.QueryStringParameters)
	if !pass {
		return validationResponse
	}

	context := context.Context{
		Body:            event.Body,
		QueryParameters: event.QueryStringParameters,
		PathParameters:  event.PathParameters,
		Headers:         event.Headers,
	}

	return handleResponse(handler, context)
}

func (r *Router) register(method http.HTTPMethod, path string, handler handler.Handler) error {
	if !method.IsValid() {
		return errors.New("invalid http method")
	}
	if len(path) == 0 {
		return errors.New("invalid resource")
	}

	if handler.Handler == nil {
		return errors.New("resource handler required")
	}

	switch method {
	case http.GET:
		r.gets[path] = handler
	case http.PATCH:
		r.patches[path] = handler
	case http.PUT:
		r.puts[path] = handler
	case http.POST:
		r.posts[path] = handler
	}
	return nil
}

func handleResponse(pathElement handler.Handler, context context.Context) events.APIGatewayProxyResponse {

	Response, err := pathElement.Handler(context)
	if err != nil {
		if pathElement.ErrorHandler == nil {
			r := defaultErrorResponse()
			if Response.Code != 0 {
				r.StatusCode = Response.Code
				if Response.Body != nil {
					jsonData, err := json.Marshal(Response.Body)
					if err != nil {
						return defaultErrorResponse()
					}
					r.Body = string(jsonData)
					return r
				}
			}
		}
		Response, err = pathElement.ErrorHandler(err)
		if err != nil {
			return defaultErrorResponse()
		}
		jsonData, err := json.Marshal(Response.Body)
		if err != nil {
			return defaultErrorResponse()
		}
		return events.APIGatewayProxyResponse{
			StatusCode: Response.Code,
			Body:       string(jsonData),
		}
	}

	jsonData, err := json.Marshal(Response.Body)
	if err != nil {
		return defaultErrorResponse()
	}
	return events.APIGatewayProxyResponse{
		StatusCode: Response.Code,
		Body:       string(jsonData),
	}
}

func defaultErrorResponse() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
	}
}

func validateQueryParameters(mandatory []string, parameters map[string]string) (events.APIGatewayProxyResponse, bool) {

	for _, e := range mandatory {
		_, ok := parameters[e]
		if !ok {
			str := ""
			for i, e := range mandatory {
				str += e
				if i != len(mandatory)-1 {
					str += ", "
				}
			}
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       fmt.Sprintf("endpoint requires the following query parameters: %s", str),
			}, false
		}
	}

	return events.APIGatewayProxyResponse{}, true
}
