package devices

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

type DownloadUpdatesOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DownloadUpdates ...
func (c DevicesClient) DownloadUpdates(ctx context.Context, id DataBoxEdgeDeviceId) (result DownloadUpdatesOperationResponse, err error) {
	req, err := c.preparerForDownloadUpdates(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "DownloadUpdates", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDownloadUpdates(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "DownloadUpdates", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DownloadUpdatesThenPoll performs DownloadUpdates then polls until it's completed
func (c DevicesClient) DownloadUpdatesThenPoll(ctx context.Context, id DataBoxEdgeDeviceId) error {
	result, err := c.DownloadUpdates(ctx, id)
	if err != nil {
		return fmt.Errorf("performing DownloadUpdates: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DownloadUpdates: %+v", err)
	}

	return nil
}

// preparerForDownloadUpdates prepares the DownloadUpdates request.
func (c DevicesClient) preparerForDownloadUpdates(ctx context.Context, id DataBoxEdgeDeviceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/downloadUpdates", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDownloadUpdates sends the DownloadUpdates request. The method will close the
// http.Response Body if it receives an error.
func (c DevicesClient) senderForDownloadUpdates(ctx context.Context, req *http.Request) (future DownloadUpdatesOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
