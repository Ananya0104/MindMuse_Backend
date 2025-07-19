package main

import (
	"context"
	"encoding/json"

	// "lambda-server/database"
	"fmt"
	"lambda-server/constants"
	"lambda-server/routes"
	"lambda-server/utils"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

// Global variables
var (
	// dynamoSvc     *database.DynamoDBService
	ginLambda     *ginadapter.GinLambdaV2
	router        *gin.Engine
)

func init() {
	// Set Gin mode
	setGinMode()

	// Setup router using the routes package
	router = routes.SetupRouter()

	// Initialize Lambda adapter if not running locally
	if !utils.IsRunningLocally() {
		ginLambda = ginadapter.NewV2(router)
	}
}

func setGinMode() {
	if os.Getenv(constants.GIN_MODE) == constants.EMPTY_STRING {
		if utils.IsRunningLocally() {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
	}
}

func Handler(ctx context.Context, event interface{}) (interface{}, error) {
	eventBytes, _ := json.Marshal(event)

	// Try to unmarshal as Lambda Function URL event
	var functionURLEvent events.LambdaFunctionURLRequest
	if err := json.Unmarshal(eventBytes, &functionURLEvent); err == nil && functionURLEvent.RequestContext.HTTP.Method != "" {
		res, err := ginLambda.ProxyWithContext(ctx, convertFunctionURLToAPIGatewayV2(functionURLEvent))
		responseBytes, _ := json.Marshal(res)
    log.Printf("Lambda Response: %s\n", string(responseBytes))
		return res, err
	}

	// Try to unmarshal as API Gateway event
	var apiGatewayEvent events.APIGatewayV2HTTPRequest
	if err := json.Unmarshal(eventBytes, &apiGatewayEvent); err == nil && apiGatewayEvent.RequestContext.HTTP.Method != "" {
		return ginLambda.ProxyWithContext(ctx, apiGatewayEvent)
	}

	// Default response for unknown event types
	return ginLambda.ProxyWithContext(ctx, events.APIGatewayV2HTTPRequest{
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "GET",
				Path:   "/health",
			},
		},
		RawPath: constants.Health,
	})
}

// convertFunctionURLToAPIGatewayV2 converts a LambdaFunctionURLRequest into an APIGatewayV2HTTPRequest.
func convertFunctionURLToAPIGatewayV2(funcURLEvent events.LambdaFunctionURLRequest) events.APIGatewayV2HTTPRequest {
	return events.APIGatewayV2HTTPRequest{
		Version:         "2.0", // API Gateway V2.0 event format version
		RouteKey:        "$default", // Often $default for Function URLs or simple API Gateway HTTP APIs
		RawPath:         funcURLEvent.RawPath,
		RawQueryString:  funcURLEvent.RawQueryString,
		Headers:         funcURLEvent.Headers,
		QueryStringParameters: funcURLEvent.QueryStringParameters,
		Body:            funcURLEvent.Body,
		IsBase64Encoded: funcURLEvent.IsBase64Encoded,
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			AccountID:    funcURLEvent.RequestContext.AccountID,
			APIID:        funcURLEvent.RequestContext.RequestID, // Placeholder, APIID is not directly from Function URL context
			DomainName:   funcURLEvent.RequestContext.DomainName,
			DomainPrefix: funcURLEvent.RequestContext.DomainPrefix,
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method:    funcURLEvent.RequestContext.HTTP.Method,
				Path:      funcURLEvent.RawPath,
				Protocol:  funcURLEvent.RequestContext.HTTP.Protocol,
				SourceIP:  funcURLEvent.RequestContext.HTTP.SourceIP,
				UserAgent: funcURLEvent.RequestContext.HTTP.UserAgent,
			},
			RequestID: funcURLEvent.RequestContext.RequestID,
			RouteKey:  "$default", // Consistent with RouteKey at top level
			Stage:     "$default", // Function URLs generally use $default stage
			Time:      funcURLEvent.RequestContext.Time,
			TimeEpoch: funcURLEvent.RequestContext.TimeEpoch,
		},
	}
}

func convertFunctionURLToAPIGateway(funcURLEvent events.LambdaFunctionURLRequest) events.APIGatewayProxyRequest {
	return events.APIGatewayProxyRequest{
		HTTPMethod:            funcURLEvent.RequestContext.HTTP.Method,
		Path:                  funcURLEvent.RawPath,
		QueryStringParameters: funcURLEvent.QueryStringParameters,
		Headers:               funcURLEvent.Headers,
		Body:                  funcURLEvent.Body,
		IsBase64Encoded:       funcURLEvent.IsBase64Encoded,
		RequestContext: events.APIGatewayProxyRequestContext{
			HTTPMethod: funcURLEvent.RequestContext.HTTP.Method,
		},
	}
}

func runLocalServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on http://localhost:%s", port)
	for _, route := range router.Routes() { 
		fmt.Printf("%s http://localhost:8080%s\n", route.Method, route.Path) 
	}

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func main() {
	if utils.IsRunningLocally() {
		runLocalServer()
	} else {
		lambda.Start(Handler)
	}
}
