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

type AssetsListContainerSasOperationResponse struct {
	HttpResponse *http.Response
	Model        *AssetContainerSas
}

// AssetsListContainerSas ...
func (c AssetsAndAssetFiltersClient) AssetsListContainerSas(ctx context.Context, id AssetId, input ListContainerSasInput) (result AssetsListContainerSasOperationResponse, err error) {
	req, err := c.preparerForAssetsListContainerSas(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsListContainerSas", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsListContainerSas", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAssetsListContainerSas(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsListContainerSas", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAssetsListContainerSas prepares the AssetsListContainerSas request.
func (c AssetsAndAssetFiltersClient) preparerForAssetsListContainerSas(ctx context.Context, id AssetId, input ListContainerSasInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listContainerSas", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAssetsListContainerSas handles the response to the AssetsListContainerSas request. The method always
// closes the http.Response Body.
func (c AssetsAndAssetFiltersClient) responderForAssetsListContainerSas(resp *http.Response) (result AssetsListContainerSasOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
