package assetsandassetfilters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetsListStreamingLocatorsOperationResponse struct {
	HttpResponse *http.Response
	Model        *ListStreamingLocatorsResponse
}

// AssetsListStreamingLocators ...
func (c AssetsAndAssetFiltersClient) AssetsListStreamingLocators(ctx context.Context, id AssetId) (result AssetsListStreamingLocatorsOperationResponse, err error) {
	req, err := c.preparerForAssetsListStreamingLocators(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsListStreamingLocators", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsListStreamingLocators", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAssetsListStreamingLocators(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsListStreamingLocators", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAssetsListStreamingLocators prepares the AssetsListStreamingLocators request.
func (c AssetsAndAssetFiltersClient) preparerForAssetsListStreamingLocators(ctx context.Context, id AssetId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listStreamingLocators", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAssetsListStreamingLocators handles the response to the AssetsListStreamingLocators request. The method always
// closes the http.Response Body.
func (c AssetsAndAssetFiltersClient) responderForAssetsListStreamingLocators(resp *http.Response) (result AssetsListStreamingLocatorsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
