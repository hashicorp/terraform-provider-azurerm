package assetsandassetfilters

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetFiltersDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// AssetFiltersDelete ...
func (c AssetsAndAssetFiltersClient) AssetFiltersDelete(ctx context.Context, id AssetFilterId) (result AssetFiltersDeleteOperationResponse, err error) {
	req, err := c.preparerForAssetFiltersDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetFiltersDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetFiltersDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAssetFiltersDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetFiltersDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAssetFiltersDelete prepares the AssetFiltersDelete request.
func (c AssetsAndAssetFiltersClient) preparerForAssetFiltersDelete(ctx context.Context, id AssetFilterId) (*http.Request, error) {
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

// responderForAssetFiltersDelete handles the response to the AssetFiltersDelete request. The method always
// closes the http.Response Body.
func (c AssetsAndAssetFiltersClient) responderForAssetFiltersDelete(resp *http.Response) (result AssetFiltersDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
