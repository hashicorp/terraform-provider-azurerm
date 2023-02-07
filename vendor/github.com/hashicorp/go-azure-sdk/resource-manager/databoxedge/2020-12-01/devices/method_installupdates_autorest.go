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

type InstallUpdatesOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// InstallUpdates ...
func (c DevicesClient) InstallUpdates(ctx context.Context, id DataBoxEdgeDeviceId) (result InstallUpdatesOperationResponse, err error) {
	req, err := c.preparerForInstallUpdates(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "InstallUpdates", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForInstallUpdates(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "InstallUpdates", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// InstallUpdatesThenPoll performs InstallUpdates then polls until it's completed
func (c DevicesClient) InstallUpdatesThenPoll(ctx context.Context, id DataBoxEdgeDeviceId) error {
	result, err := c.InstallUpdates(ctx, id)
	if err != nil {
		return fmt.Errorf("performing InstallUpdates: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after InstallUpdates: %+v", err)
	}

	return nil
}

// preparerForInstallUpdates prepares the InstallUpdates request.
func (c DevicesClient) preparerForInstallUpdates(ctx context.Context, id DataBoxEdgeDeviceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/installUpdates", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForInstallUpdates sends the InstallUpdates request. The method will close the
// http.Response Body if it receives an error.
func (c DevicesClient) senderForInstallUpdates(ctx context.Context, req *http.Request) (future InstallUpdatesOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
