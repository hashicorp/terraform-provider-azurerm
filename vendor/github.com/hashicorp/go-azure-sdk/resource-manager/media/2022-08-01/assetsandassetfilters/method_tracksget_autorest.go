package assetsandassetfilters

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TracksGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *AssetTrack
}

// TracksGet ...
func (c AssetsAndAssetFiltersClient) TracksGet(ctx context.Context, id TrackId) (result TracksGetOperationResponse, err error) {
	req, err := c.preparerForTracksGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "TracksGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "TracksGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTracksGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "TracksGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTracksGet prepares the TracksGet request.
func (c AssetsAndAssetFiltersClient) preparerForTracksGet(ctx context.Context, id TrackId) (*http.Request, error) {
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

// responderForTracksGet handles the response to the TracksGet request. The method always
// closes the http.Response Body.
func (c AssetsAndAssetFiltersClient) responderForTracksGet(resp *http.Response) (result TracksGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
