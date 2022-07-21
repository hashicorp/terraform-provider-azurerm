package resource

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SpatialAnchorsAccountsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// SpatialAnchorsAccountsDelete ...
func (c ResourceClient) SpatialAnchorsAccountsDelete(ctx context.Context, id SpatialAnchorsAccountId) (result SpatialAnchorsAccountsDeleteOperationResponse, err error) {
	req, err := c.preparerForSpatialAnchorsAccountsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSpatialAnchorsAccountsDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSpatialAnchorsAccountsDelete prepares the SpatialAnchorsAccountsDelete request.
func (c ResourceClient) preparerForSpatialAnchorsAccountsDelete(ctx context.Context, id SpatialAnchorsAccountId) (*http.Request, error) {
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

// responderForSpatialAnchorsAccountsDelete handles the response to the SpatialAnchorsAccountsDelete request. The method always
// closes the http.Response Body.
func (c ResourceClient) responderForSpatialAnchorsAccountsDelete(resp *http.Response) (result SpatialAnchorsAccountsDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
