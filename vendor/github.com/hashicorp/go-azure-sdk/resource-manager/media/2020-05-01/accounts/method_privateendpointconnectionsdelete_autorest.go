package accounts

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointConnectionsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// PrivateEndpointConnectionsDelete ...
func (c AccountsClient) PrivateEndpointConnectionsDelete(ctx context.Context, id PrivateEndpointConnectionId) (result PrivateEndpointConnectionsDeleteOperationResponse, err error) {
	req, err := c.preparerForPrivateEndpointConnectionsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "PrivateEndpointConnectionsDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "PrivateEndpointConnectionsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPrivateEndpointConnectionsDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "PrivateEndpointConnectionsDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPrivateEndpointConnectionsDelete prepares the PrivateEndpointConnectionsDelete request.
func (c AccountsClient) preparerForPrivateEndpointConnectionsDelete(ctx context.Context, id PrivateEndpointConnectionId) (*http.Request, error) {
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

// responderForPrivateEndpointConnectionsDelete handles the response to the PrivateEndpointConnectionsDelete request. The method always
// closes the http.Response Body.
func (c AccountsClient) responderForPrivateEndpointConnectionsDelete(resp *http.Response) (result PrivateEndpointConnectionsDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
