package resource

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SpatialAnchorsAccountsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *SpatialAnchorsAccount
}

// SpatialAnchorsAccountsGet ...
func (c ResourceClient) SpatialAnchorsAccountsGet(ctx context.Context, id SpatialAnchorsAccountId) (result SpatialAnchorsAccountsGetOperationResponse, err error) {
	req, err := c.preparerForSpatialAnchorsAccountsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSpatialAnchorsAccountsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSpatialAnchorsAccountsGet prepares the SpatialAnchorsAccountsGet request.
func (c ResourceClient) preparerForSpatialAnchorsAccountsGet(ctx context.Context, id SpatialAnchorsAccountId) (*http.Request, error) {
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

// responderForSpatialAnchorsAccountsGet handles the response to the SpatialAnchorsAccountsGet request. The method always
// closes the http.Response Body.
func (c ResourceClient) responderForSpatialAnchorsAccountsGet(resp *http.Response) (result SpatialAnchorsAccountsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
