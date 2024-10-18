package router

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type Router struct {
	gets    map[string]PathElement
	puts    map[string]PathElement
	patches map[string]PathElement
	posts   map[string]PathElement
}

type PathElement struct {
	context                  RouteContext
	handler                  func(RouteContext) (HttpResponse, error)
	errorHandler             func(error) (HttpResponse, error)
	mandatoryQueryParameters []string
}

type HttpResponse struct {
	Code int
	Body any
}

type RouteContext struct {
	PathParameters  map[string]string
	QueryParameters map[string]string
	Headers         map[string]string
}

func New() *Router {
	return &Router{
		gets: map[string]PathElement{},
	}
}

func (r *Router) RegisterGet(
	path string,
	handler func(context RouteContext) (HttpResponse, error),
	errorHandler func(err error) (HttpResponse, error),
	mandatoryQueryParameters []string) {

	r.register(GET, path, handler, errorHandler, mandatoryQueryParameters)
}

func (r *Router) RegisterPut(
	path string,
	handler func(context RouteContext) (HttpResponse, error),
	errorHandler func(err error) (HttpResponse, error),
	mandatoryQueryParameters []string) {

	r.register(PUT, path, handler, errorHandler, mandatoryQueryParameters)
}

func (r *Router) RegisterPatch(
	path string,
	handler func(context RouteContext) (HttpResponse, error),
	errorHandler func(err error) (HttpResponse, error),
	mandatoryQueryParameters []string) {

	r.register(PATCH, path, handler, errorHandler, mandatoryQueryParameters)
}

func (r *Router) RegisterPost(
	path string,
	handler func(context RouteContext) (HttpResponse, error),
	errorHandler func(err error) (HttpResponse, error),
	mandatoryQueryParameters []string) {

	r.register(POST, path, handler, errorHandler, mandatoryQueryParameters)
}

func (r *Router) Route(event events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	httpMethod, err := HttpMethodFrom(event.HTTPMethod)

	if err != nil {
		return defaultErrorResponse()
	}

	switch httpMethod {
	case GET:
		return route(r.gets, event)
	case PUT:
		return route(r.puts, event)
	case PATCH:
		return route(r.patches, event)
	case POST:
		return route(r.posts, event)
	default:
		response := events.APIGatewayProxyResponse{
			StatusCode: 405,
		}
		return response
	}

}
func route(methodPaths map[string]PathElement, event events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	pathElement, ok := methodPaths[event.Resource]
	if !ok {
		response := events.APIGatewayProxyResponse{
			StatusCode: 404,
		}
		return response
	}
	pathElement.context.QueryParameters = event.QueryStringParameters
	pathElement.context.PathParameters = event.PathParameters
	validationResponse, pass := validateQueryParameters(pathElement.mandatoryQueryParameters, event.QueryStringParameters)
	if !pass {
		return validationResponse
	}
	return handleResponse(pathElement)
}

func (r *Router) register(
	method HTTPMethod,
	path string,
	getFunc func(context RouteContext) (HttpResponse, error),
	errorHandler func(err error) (HttpResponse, error),
	mandatoryQueryParameters []string) {

	element := PathElement{
		handler:                  getFunc,
		errorHandler:             errorHandler,
		mandatoryQueryParameters: mandatoryQueryParameters,
	}

	switch method {
	case GET:
		r.gets[path] = element
	}
}

func handleResponse(pathElement PathElement) events.APIGatewayProxyResponse {

	httpResponse, err := pathElement.handler(pathElement.context)
	if err != nil {
		if pathElement.errorHandler == nil {
			r := defaultErrorResponse()
			if httpResponse.Code != 0 {
				r.StatusCode = httpResponse.Code
				if httpResponse.Body != nil {
					jsonData, err := json.Marshal(httpResponse.Body)
					if err != nil {
						return defaultErrorResponse()
					}
					r.Body = string(jsonData)
					return r
				}
			}
		}
		httpResponse, err = pathElement.errorHandler(err)
		if err != nil {
			return defaultErrorResponse()
		}
		jsonData, err := json.Marshal(httpResponse.Body)
		if err != nil {
			return defaultErrorResponse()
		}
		return events.APIGatewayProxyResponse{
			StatusCode: httpResponse.Code,
			Body:       string(jsonData),
		}
	}

	jsonData, err := json.Marshal(httpResponse.Body)
	if err != nil {
		return defaultErrorResponse()
	}
	return events.APIGatewayProxyResponse{
		StatusCode: httpResponse.Code,
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
