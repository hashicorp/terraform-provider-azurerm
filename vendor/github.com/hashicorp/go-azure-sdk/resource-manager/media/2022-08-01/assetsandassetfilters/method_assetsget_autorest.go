package assetsandassetfilters

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *Asset
}

// AssetsGet ...
func (c AssetsAndAssetFiltersClient) AssetsGet(ctx context.Context, id AssetId) (result AssetsGetOperationResponse, err error) {
	req, err := c.preparerForAssetsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAssetsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAssetsGet prepares the AssetsGet request.
func (c AssetsAndAssetFiltersClient) preparerForAssetsGet(ctx context.Context, id AssetId) (*http.Request, error) {
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

// responderForAssetsGet handles the response to the AssetsGet request. The method always
// closes the http.Response Body.
func (c AssetsAndAssetFiltersClient) responderForAssetsGet(resp *http.Response) (result AssetsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
