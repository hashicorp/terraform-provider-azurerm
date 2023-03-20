package redis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ForceRebootOperationResponse struct {
	HttpResponse *http.Response
	Model        *RedisForceRebootResponse
}

// ForceReboot ...
func (c RedisClient) ForceReboot(ctx context.Context, id RediId, input RedisRebootParameters) (result ForceRebootOperationResponse, err error) {
	req, err := c.preparerForForceReboot(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "ForceReboot", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "ForceReboot", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForForceReboot(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "ForceReboot", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForForceReboot prepares the ForceReboot request.
func (c RedisClient) preparerForForceReboot(ctx context.Context, id RediId, input RedisRebootParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/forceReboot", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForForceReboot handles the response to the ForceReboot request. The method always
// closes the http.Response Body.
func (c RedisClient) responderForForceReboot(resp *http.Response) (result ForceRebootOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
