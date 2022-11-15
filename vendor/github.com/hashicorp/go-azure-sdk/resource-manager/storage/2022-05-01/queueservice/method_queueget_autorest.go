package queueservice

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueueGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *StorageQueue
}

// QueueGet ...
func (c QueueServiceClient) QueueGet(ctx context.Context, id QueueId) (result QueueGetOperationResponse, err error) {
	req, err := c.preparerForQueueGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueservice.QueueServiceClient", "QueueGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueservice.QueueServiceClient", "QueueGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForQueueGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueservice.QueueServiceClient", "QueueGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForQueueGet prepares the QueueGet request.
func (c QueueServiceClient) preparerForQueueGet(ctx context.Context, id QueueId) (*http.Request, error) {
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

// responderForQueueGet handles the response to the QueueGet request. The method always
// closes the http.Response Body.
func (c QueueServiceClient) responderForQueueGet(resp *http.Response) (result QueueGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
