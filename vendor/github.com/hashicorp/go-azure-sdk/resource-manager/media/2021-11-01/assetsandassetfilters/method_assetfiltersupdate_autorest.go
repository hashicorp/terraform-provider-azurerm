package assetsandassetfilters

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetFiltersUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *AssetFilter
}

// AssetFiltersUpdate ...
func (c AssetsAndAssetFiltersClient) AssetFiltersUpdate(ctx context.Context, id AssetFilterId, input AssetFilter) (result AssetFiltersUpdateOperationResponse, err error) {
	req, err := c.preparerForAssetFiltersUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetFiltersUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetFiltersUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAssetFiltersUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetFiltersUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAssetFiltersUpdate prepares the AssetFiltersUpdate request.
func (c AssetsAndAssetFiltersClient) preparerForAssetFiltersUpdate(ctx context.Context, id AssetFilterId, input AssetFilter) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAssetFiltersUpdate handles the response to the AssetFiltersUpdate request. The method always
// closes the http.Response Body.
func (c AssetsAndAssetFiltersClient) responderForAssetFiltersUpdate(resp *http.Response) (result AssetFiltersUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
