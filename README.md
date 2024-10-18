# Background
This module is intended to route AWS API Gateway proxy events for lambda functions.

Users are able to register handlers for each API resource along with an error handler and define mandatory query parameters.

If the client fails to pass required query parmeters the router will respond with the following:
```json
{
    "StatusCode" : 400,
    "Body" : "endpoint requires the following query paramters: p1, p2"
}
```

If the client attempts to hit and resource which is not registered, the following will be returned:
```json
{
    "StatusCode" : 404
}
```


# Example Code

```go

func Handler(_ context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // This is only an example! The router generally should be initialized once and reused!
    router := router.New()
    router.RegisterGet(
        "/v1/some/path"
        handler,
        errorHandler,
        [] string {"p1"}
    )

    return router.Route(r), nil
}

```