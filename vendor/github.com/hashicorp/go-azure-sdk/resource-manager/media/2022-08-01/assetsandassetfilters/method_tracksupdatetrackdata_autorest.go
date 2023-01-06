package assetsandassetfilters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TracksUpdateTrackDataOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// TracksUpdateTrackData ...
func (c AssetsAndAssetFiltersClient) TracksUpdateTrackData(ctx context.Context, id TrackId) (result TracksUpdateTrackDataOperationResponse, err error) {
	req, err := c.preparerForTracksUpdateTrackData(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "TracksUpdateTrackData", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForTracksUpdateTrackData(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "TracksUpdateTrackData", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// TracksUpdateTrackDataThenPoll performs TracksUpdateTrackData then polls until it's completed
func (c AssetsAndAssetFiltersClient) TracksUpdateTrackDataThenPoll(ctx context.Context, id TrackId) error {
	result, err := c.TracksUpdateTrackData(ctx, id)
	if err != nil {
		return fmt.Errorf("performing TracksUpdateTrackData: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after TracksUpdateTrackData: %+v", err)
	}

	return nil
}

// preparerForTracksUpdateTrackData prepares the TracksUpdateTrackData request.
func (c AssetsAndAssetFiltersClient) preparerForTracksUpdateTrackData(ctx context.Context, id TrackId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/updateTrackData", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForTracksUpdateTrackData sends the TracksUpdateTrackData request. The method will close the
// http.Response Body if it receives an error.
func (c AssetsAndAssetFiltersClient) senderForTracksUpdateTrackData(ctx context.Context, req *http.Request) (future TracksUpdateTrackDataOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
