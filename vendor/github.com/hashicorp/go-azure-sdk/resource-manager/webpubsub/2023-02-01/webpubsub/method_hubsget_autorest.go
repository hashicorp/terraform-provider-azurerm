package webpubsub

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HubsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *WebPubSubHub
}

// HubsGet ...
func (c WebPubSubClient) HubsGet(ctx context.Context, id HubId) (result HubsGetOperationResponse, err error) {
	req, err := c.preparerForHubsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webpubsub.WebPubSubClient", "HubsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "webpubsub.WebPubSubClient", "HubsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForHubsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webpubsub.WebPubSubClient", "HubsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForHubsGet prepares the HubsGet request.
func (c WebPubSubClient) preparerForHubsGet(ctx context.Context, id HubId) (*http.Request, error) {
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

// responderForHubsGet handles the response to the HubsGet request. The method always
// closes the http.Response Body.
func (c WebPubSubClient) responderForHubsGet(resp *http.Response) (result HubsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
