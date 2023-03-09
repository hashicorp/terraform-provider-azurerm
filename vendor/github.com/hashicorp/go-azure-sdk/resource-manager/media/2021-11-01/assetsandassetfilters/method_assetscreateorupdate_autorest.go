package assetsandassetfilters

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetsCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Asset
}

// AssetsCreateOrUpdate ...
func (c AssetsAndAssetFiltersClient) AssetsCreateOrUpdate(ctx context.Context, id AssetId, input Asset) (result AssetsCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForAssetsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAssetsCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAssetsCreateOrUpdate prepares the AssetsCreateOrUpdate request.
func (c AssetsAndAssetFiltersClient) preparerForAssetsCreateOrUpdate(ctx context.Context, id AssetId, input Asset) (*http.Request, error) {
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

// responderForAssetsCreateOrUpdate handles the response to the AssetsCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c AssetsAndAssetFiltersClient) responderForAssetsCreateOrUpdate(resp *http.Response) (result AssetsCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
