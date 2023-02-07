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

type TracksCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// TracksCreateOrUpdate ...
func (c AssetsAndAssetFiltersClient) TracksCreateOrUpdate(ctx context.Context, id TrackId, input AssetTrack) (result TracksCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForTracksCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "TracksCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForTracksCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "TracksCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// TracksCreateOrUpdateThenPoll performs TracksCreateOrUpdate then polls until it's completed
func (c AssetsAndAssetFiltersClient) TracksCreateOrUpdateThenPoll(ctx context.Context, id TrackId, input AssetTrack) error {
	result, err := c.TracksCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing TracksCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after TracksCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForTracksCreateOrUpdate prepares the TracksCreateOrUpdate request.
func (c AssetsAndAssetFiltersClient) preparerForTracksCreateOrUpdate(ctx context.Context, id TrackId, input AssetTrack) (*http.Request, error) {
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

// senderForTracksCreateOrUpdate sends the TracksCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c AssetsAndAssetFiltersClient) senderForTracksCreateOrUpdate(ctx context.Context, req *http.Request) (future TracksCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
