package assetsandassetfilters

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// AssetsDelete ...
func (c AssetsAndAssetFiltersClient) AssetsDelete(ctx context.Context, id AssetId) (result AssetsDeleteOperationResponse, err error) {
	req, err := c.preparerForAssetsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAssetsDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAssetsDelete prepares the AssetsDelete request.
func (c AssetsAndAssetFiltersClient) preparerForAssetsDelete(ctx context.Context, id AssetId) (*http.Request, error) {
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

// responderForAssetsDelete handles the response to the AssetsDelete request. The method always
// closes the http.Response Body.
func (c AssetsAndAssetFiltersClient) responderForAssetsDelete(resp *http.Response) (result AssetsDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
