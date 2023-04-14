package redis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkedServerGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *RedisLinkedServerWithProperties
}

// LinkedServerGet ...
func (c RedisClient) LinkedServerGet(ctx context.Context, id LinkedServerId) (result LinkedServerGetOperationResponse, err error) {
	req, err := c.preparerForLinkedServerGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "LinkedServerGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "LinkedServerGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForLinkedServerGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "LinkedServerGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForLinkedServerGet prepares the LinkedServerGet request.
func (c RedisClient) preparerForLinkedServerGet(ctx context.Context, id LinkedServerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForLinkedServerGet handles the response to the LinkedServerGet request. The method always
// closes the http.Response Body.
func (c RedisClient) responderForLinkedServerGet(resp *http.Response) (result LinkedServerGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
