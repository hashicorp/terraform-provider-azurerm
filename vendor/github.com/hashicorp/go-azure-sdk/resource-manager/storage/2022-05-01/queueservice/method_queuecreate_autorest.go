package queueservice

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueueCreateOperationResponse struct {
	HttpResponse *http.Response
	Model        *StorageQueue
}

// QueueCreate ...
func (c QueueServiceClient) QueueCreate(ctx context.Context, id QueueId, input StorageQueue) (result QueueCreateOperationResponse, err error) {
	req, err := c.preparerForQueueCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueservice.QueueServiceClient", "QueueCreate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueservice.QueueServiceClient", "QueueCreate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForQueueCreate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queueservice.QueueServiceClient", "QueueCreate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForQueueCreate prepares the QueueCreate request.
func (c QueueServiceClient) preparerForQueueCreate(ctx context.Context, id QueueId, input StorageQueue) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForQueueCreate handles the response to the QueueCreate request. The method always
// closes the http.Response Body.
func (c QueueServiceClient) responderForQueueCreate(resp *http.Response) (result QueueCreateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
