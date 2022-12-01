package assetsandassetfilters

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetsUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Asset
}

// AssetsUpdate ...
func (c AssetsAndAssetFiltersClient) AssetsUpdate(ctx context.Context, id AssetId, input Asset) (result AssetsUpdateOperationResponse, err error) {
	req, err := c.preparerForAssetsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAssetsUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAssetsUpdate prepares the AssetsUpdate request.
func (c AssetsAndAssetFiltersClient) preparerForAssetsUpdate(ctx context.Context, id AssetId, input Asset) (*http.Request, error) {
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

// responderForAssetsUpdate handles the response to the AssetsUpdate request. The method always
// closes the http.Response Body.
func (c AssetsAndAssetFiltersClient) responderForAssetsUpdate(resp *http.Response) (result AssetsUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
