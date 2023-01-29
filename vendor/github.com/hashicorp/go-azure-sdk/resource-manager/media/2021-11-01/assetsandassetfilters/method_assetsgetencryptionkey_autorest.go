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

type AssetsGetEncryptionKeyOperationResponse struct {
	HttpResponse *http.Response
	Model        *StorageEncryptedAssetDecryptionData
}

// AssetsGetEncryptionKey ...
func (c AssetsAndAssetFiltersClient) AssetsGetEncryptionKey(ctx context.Context, id AssetId) (result AssetsGetEncryptionKeyOperationResponse, err error) {
	req, err := c.preparerForAssetsGetEncryptionKey(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsGetEncryptionKey", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsGetEncryptionKey", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAssetsGetEncryptionKey(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsGetEncryptionKey", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAssetsGetEncryptionKey prepares the AssetsGetEncryptionKey request.
func (c AssetsAndAssetFiltersClient) preparerForAssetsGetEncryptionKey(ctx context.Context, id AssetId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/getEncryptionKey", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAssetsGetEncryptionKey handles the response to the AssetsGetEncryptionKey request. The method always
// closes the http.Response Body.
func (c AssetsAndAssetFiltersClient) responderForAssetsGetEncryptionKey(resp *http.Response) (result AssetsGetEncryptionKeyOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
