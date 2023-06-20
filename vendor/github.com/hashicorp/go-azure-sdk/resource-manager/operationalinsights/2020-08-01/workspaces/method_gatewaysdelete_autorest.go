package workspaces

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewaysDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// GatewaysDelete ...
func (c WorkspacesClient) GatewaysDelete(ctx context.Context, id GatewayId) (result GatewaysDeleteOperationResponse, err error) {
	req, err := c.preparerForGatewaysDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "GatewaysDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "GatewaysDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGatewaysDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "GatewaysDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGatewaysDelete prepares the GatewaysDelete request.
func (c WorkspacesClient) preparerForGatewaysDelete(ctx context.Context, id GatewayId) (*http.Request, error) {
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

// responderForGatewaysDelete handles the response to the GatewaysDelete request. The method always
// closes the http.Response Body.
func (c WorkspacesClient) responderForGatewaysDelete(resp *http.Response) (result GatewaysDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
