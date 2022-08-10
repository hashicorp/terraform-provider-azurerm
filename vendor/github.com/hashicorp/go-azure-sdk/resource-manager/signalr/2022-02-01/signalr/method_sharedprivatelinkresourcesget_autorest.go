package signalr

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SharedPrivateLinkResourcesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *SharedPrivateLinkResource
}

// SharedPrivateLinkResourcesGet ...
func (c SignalRClient) SharedPrivateLinkResourcesGet(ctx context.Context, id SharedPrivateLinkResourceId) (result SharedPrivateLinkResourcesGetOperationResponse, err error) {
	req, err := c.preparerForSharedPrivateLinkResourcesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "SharedPrivateLinkResourcesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "SharedPrivateLinkResourcesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSharedPrivateLinkResourcesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "SharedPrivateLinkResourcesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSharedPrivateLinkResourcesGet prepares the SharedPrivateLinkResourcesGet request.
func (c SignalRClient) preparerForSharedPrivateLinkResourcesGet(ctx context.Context, id SharedPrivateLinkResourceId) (*http.Request, error) {
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

// responderForSharedPrivateLinkResourcesGet handles the response to the SharedPrivateLinkResourcesGet request. The method always
// closes the http.Response Body.
func (c SignalRClient) responderForSharedPrivateLinkResourcesGet(resp *http.Response) (result SharedPrivateLinkResourcesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
