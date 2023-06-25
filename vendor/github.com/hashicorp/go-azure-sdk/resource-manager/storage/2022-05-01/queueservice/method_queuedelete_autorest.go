package queueservice

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueueDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// QueueDelete ...
func (c QueueServiceClient) QueueDelete(ctx context.Context, id QueueId) (result QueueDeleteOperationResponse, err error) {
	req, err := c.preparerForQueueDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueservice.QueueServiceClient", "QueueDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueservice.QueueServiceClient", "QueueDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForQueueDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueservice.QueueServiceClient", "QueueDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForQueueDelete prepares the QueueDelete request.
func (c QueueServiceClient) preparerForQueueDelete(ctx context.Context, id QueueId) (*http.Request, error) {
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

// responderForQueueDelete handles the response to the QueueDelete request. The method always
// closes the http.Response Body.
func (c QueueServiceClient) responderForQueueDelete(resp *http.Response) (result QueueDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
