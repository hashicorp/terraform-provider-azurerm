package redis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkedServerDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// LinkedServerDelete ...
func (c RedisClient) LinkedServerDelete(ctx context.Context, id LinkedServerId) (result LinkedServerDeleteOperationResponse, err error) {
	req, err := c.preparerForLinkedServerDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "LinkedServerDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "LinkedServerDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForLinkedServerDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "LinkedServerDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForLinkedServerDelete prepares the LinkedServerDelete request.
func (c RedisClient) preparerForLinkedServerDelete(ctx context.Context, id LinkedServerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForLinkedServerDelete handles the response to the LinkedServerDelete request. The method always
// closes the http.Response Body.
func (c RedisClient) responderForLinkedServerDelete(resp *http.Response) (result LinkedServerDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
