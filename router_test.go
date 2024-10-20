package router

import (
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/beglaryh/aws-lambda-router/context"
	"github.com/beglaryh/aws-lambda-router/handler"
	"github.com/beglaryh/aws-lambda-router/http"
)

func TestSimpleGet(t *testing.T) {
	router := New()

	router.RegisterGet(
		"/path",
		handler.Builder().Handler(getFunction).ErrorHandler(errorFunction).Build(),
	)
	event := events.APIGatewayProxyRequest{
		HTTPMethod:            "GET",
		Resource:              "/path",
		QueryStringParameters: map[string]string{"param1": "value"},
	}
	response := router.Route(event)
	if response.StatusCode != 200 {
		t.Fail()
	}

	body := response.Body
	expectedBody := "{\"Hello\":\"World!\"}"
	if body != expectedBody {
		t.Fail()
	}

	event.QueryStringParameters = map[string]string{}
	response = router.Route(event)
	if response.StatusCode != 400 {
		t.Fail()
	}
}

func TestMandatoryQueryParameters(t *testing.T) {
	router := New()
	router.RegisterGet(
		"/path",
		handler.Builder().
			Handler(getFunction).
			ErrorHandler(errorFunction).
			MandatoryQueryParameters([]string{"param1"}).
			Build(),
	)
	event := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Resource:   "/path",
	}
	response := router.Route(event)
	if response.StatusCode != 400 {
		t.Fail()
	}
	if response.Body != "endpoint requires the following query parameters: param1" {
		t.Fail()
	}
}

func getFunction(context context.Context) (http.Response, error) {
	_, ok := context.QueryParameters["param1"]
	if !ok {
		return http.Response{Code: 400}, errors.New("Missing parameter")
	}
	return http.Response{Code: 200, Body: map[string]string{"Hello": "World!"}}, nil
}

func errorFunction(err error) (http.Response, error) {
	return http.Response{Code: 400, Body: err}, nil
}
