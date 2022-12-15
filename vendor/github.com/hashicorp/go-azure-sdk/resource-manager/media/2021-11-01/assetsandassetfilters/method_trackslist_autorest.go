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

type TracksListOperationResponse struct {
	HttpResponse *http.Response
	Model        *AssetTrackCollection
}

// TracksList ...
func (c AssetsAndAssetFiltersClient) TracksList(ctx context.Context, id AssetId) (result TracksListOperationResponse, err error) {
	req, err := c.preparerForTracksList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "TracksList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "TracksList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTracksList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "TracksList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTracksList prepares the TracksList request.
func (c AssetsAndAssetFiltersClient) preparerForTracksList(ctx context.Context, id AssetId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/tracks", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForTracksList handles the response to the TracksList request. The method always
// closes the http.Response Body.
func (c AssetsAndAssetFiltersClient) responderForTracksList(resp *http.Response) (result TracksListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
