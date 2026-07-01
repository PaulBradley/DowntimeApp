package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func main() {
	log.SetFlags(0)
	logJSON("INFO", "cold_start", map[string]any{
		"service": "downtimeapp-frontend",
	})
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (resp events.LambdaFunctionURLResponse, err error) {
	startedAt := time.Now()
	requestID := ""
	if lc, ok := lambdacontext.FromContext(ctx); ok {
		requestID = lc.AwsRequestID
	}

	logJSON("INFO", "request_received", map[string]any{
		"request_id": requestID,
		"method":     request.RequestContext.HTTP.Method,
		"path":       request.RawPath,
		"source_ip":  request.RequestContext.HTTP.SourceIP,
		"user_agent": request.RequestContext.HTTP.UserAgent,
	})

	defer func() {
		durationMs := time.Since(startedAt).Milliseconds()

		if r := recover(); r != nil {
			logJSON("ERROR", "handler_panic", map[string]any{
				"request_id":  requestID,
				"panic":       fmt.Sprint(r),
				"stack_trace": string(debug.Stack()),
				"duration_ms": durationMs,
			})
			resp = Response(500, "internal server error")
			err = nil
			return
		}

		statusCode := resp.StatusCode
		if statusCode == 0 {
			statusCode = 500
		}

		logJSON("INFO", "request_complete", map[string]any{
			"request_id":  requestID,
			"status_code": statusCode,
			"duration_ms": durationMs,
		})
	}()

	body := "THE DOWNTIMEAPP.CLOUD FRONTEND"
	resp = Response(200, body)
	return resp, nil
}

func Response(StatusCode int, Body string) events.LambdaFunctionURLResponse {
	return events.LambdaFunctionURLResponse{
		IsBase64Encoded: false,
		Body:            Body,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
		StatusCode: StatusCode,
	}
}

func logJSON(level, message string, fields map[string]any) {
	entry := map[string]any{
		"level":   level,
		"message": message,
		"time":    time.Now().UTC().Format(time.RFC3339Nano),
	}

	for k, v := range fields {
		entry[k] = v
	}

	b, err := json.Marshal(entry)
	if err != nil {
		log.Printf("{\"level\":\"ERROR\",\"message\":\"log_marshal_failed\",\"error\":%q}", err.Error())
		return
	}

	log.Print(string(b))
}
