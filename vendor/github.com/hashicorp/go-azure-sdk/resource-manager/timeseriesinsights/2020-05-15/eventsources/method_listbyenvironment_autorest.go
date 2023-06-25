package eventsources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByEnvironmentOperationResponse struct {
	HttpResponse *http.Response
	Model        *EventSourceListResponse
}

// ListByEnvironment ...
func (c EventSourcesClient) ListByEnvironment(ctx context.Context, id EnvironmentId) (result ListByEnvironmentOperationResponse, err error) {
	req, err := c.preparerForListByEnvironment(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsources.EventSourcesClient", "ListByEnvironment", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsources.EventSourcesClient", "ListByEnvironment", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListByEnvironment(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsources.EventSourcesClient", "ListByEnvironment", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListByEnvironment prepares the ListByEnvironment request.
func (c EventSourcesClient) preparerForListByEnvironment(ctx context.Context, id EnvironmentId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/eventSources", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListByEnvironment handles the response to the ListByEnvironment request. The method always
// closes the http.Response Body.
func (c EventSourcesClient) responderForListByEnvironment(resp *http.Response) (result ListByEnvironmentOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
