package assetsandassetfilters

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetFiltersCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *AssetFilter
}

// AssetFiltersCreateOrUpdate ...
func (c AssetsAndAssetFiltersClient) AssetFiltersCreateOrUpdate(ctx context.Context, id AssetFilterId, input AssetFilter) (result AssetFiltersCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForAssetFiltersCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetFiltersCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetFiltersCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAssetFiltersCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetFiltersCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAssetFiltersCreateOrUpdate prepares the AssetFiltersCreateOrUpdate request.
func (c AssetsAndAssetFiltersClient) preparerForAssetFiltersCreateOrUpdate(ctx context.Context, id AssetFilterId, input AssetFilter) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAssetFiltersCreateOrUpdate handles the response to the AssetFiltersCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c AssetsAndAssetFiltersClient) responderForAssetFiltersCreateOrUpdate(resp *http.Response) (result AssetFiltersCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
