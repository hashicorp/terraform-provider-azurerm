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

type ScanForUpdatesOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ScanForUpdates ...
func (c DevicesClient) ScanForUpdates(ctx context.Context, id DataBoxEdgeDeviceId) (result ScanForUpdatesOperationResponse, err error) {
	req, err := c.preparerForScanForUpdates(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "ScanForUpdates", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForScanForUpdates(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "ScanForUpdates", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ScanForUpdatesThenPoll performs ScanForUpdates then polls until it's completed
func (c DevicesClient) ScanForUpdatesThenPoll(ctx context.Context, id DataBoxEdgeDeviceId) error {
	result, err := c.ScanForUpdates(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ScanForUpdates: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ScanForUpdates: %+v", err)
	}

	return nil
}

// preparerForScanForUpdates prepares the ScanForUpdates request.
func (c DevicesClient) preparerForScanForUpdates(ctx context.Context, id DataBoxEdgeDeviceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/scanForUpdates", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForScanForUpdates sends the ScanForUpdates request. The method will close the
// http.Response Body if it receives an error.
func (c DevicesClient) senderForScanForUpdates(ctx context.Context, req *http.Request) (future ScanForUpdatesOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
