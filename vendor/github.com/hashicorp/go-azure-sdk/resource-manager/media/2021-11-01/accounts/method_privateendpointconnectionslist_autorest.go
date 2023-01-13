package accounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointConnectionsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *PrivateEndpointConnectionListResult
}

// PrivateEndpointConnectionsList ...
func (c AccountsClient) PrivateEndpointConnectionsList(ctx context.Context, id MediaServiceId) (result PrivateEndpointConnectionsListOperationResponse, err error) {
	req, err := c.preparerForPrivateEndpointConnectionsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "PrivateEndpointConnectionsList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "PrivateEndpointConnectionsList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPrivateEndpointConnectionsList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "PrivateEndpointConnectionsList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPrivateEndpointConnectionsList prepares the PrivateEndpointConnectionsList request.
func (c AccountsClient) preparerForPrivateEndpointConnectionsList(ctx context.Context, id MediaServiceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/privateEndpointConnections", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForPrivateEndpointConnectionsList handles the response to the PrivateEndpointConnectionsList request. The method always
// closes the http.Response Body.
func (c AccountsClient) responderForPrivateEndpointConnectionsList(resp *http.Response) (result PrivateEndpointConnectionsListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
