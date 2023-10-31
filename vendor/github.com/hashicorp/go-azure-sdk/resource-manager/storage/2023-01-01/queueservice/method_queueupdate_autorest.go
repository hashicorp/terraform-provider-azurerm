package queueservice

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueueUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *StorageQueue
}

// QueueUpdate ...
func (c QueueServiceClient) QueueUpdate(ctx context.Context, id QueueId, input StorageQueue) (result QueueUpdateOperationResponse, err error) {
	req, err := c.preparerForQueueUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueservice.QueueServiceClient", "QueueUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueservice.QueueServiceClient", "QueueUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForQueueUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueservice.QueueServiceClient", "QueueUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForQueueUpdate prepares the QueueUpdate request.
func (c QueueServiceClient) preparerForQueueUpdate(ctx context.Context, id QueueId, input StorageQueue) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForQueueUpdate handles the response to the QueueUpdate request. The method always
// closes the http.Response Body.
func (c QueueServiceClient) responderForQueueUpdate(resp *http.Response) (result QueueUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
