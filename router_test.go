package router

import (
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestSimpleGet(t *testing.T) {
	router := New()
	router.RegisterGet(
		"/path",
		getFunction,
		errorFunction,
		nil,
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
		getFunction,
		errorFunction,
		[]string{"param1"},
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

func getFunction(context RouteContext) (HttpResponse, error) {
	_, ok := context.QueryParameters["param1"]
	if !ok {
		return HttpResponse{Code: 400}, errors.New("Missing parameter")
	}
	return HttpResponse{Code: 200, Body: map[string]string{"Hello": "World!"}}, nil
}

func errorFunction(err error) (HttpResponse, error) {
	return HttpResponse{Code: 400, Body: err}, nil
}
