package assetsandassetfilters

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetFiltersGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *AssetFilter
}

// AssetFiltersGet ...
func (c AssetsAndAssetFiltersClient) AssetFiltersGet(ctx context.Context, id AssetFilterId) (result AssetFiltersGetOperationResponse, err error) {
	req, err := c.preparerForAssetFiltersGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetFiltersGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetFiltersGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAssetFiltersGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetFiltersGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAssetFiltersGet prepares the AssetFiltersGet request.
func (c AssetsAndAssetFiltersClient) preparerForAssetFiltersGet(ctx context.Context, id AssetFilterId) (*http.Request, error) {
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

// responderForAssetFiltersGet handles the response to the AssetFiltersGet request. The method always
// closes the http.Response Body.
func (c AssetsAndAssetFiltersClient) responderForAssetFiltersGet(resp *http.Response) (result AssetFiltersGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
