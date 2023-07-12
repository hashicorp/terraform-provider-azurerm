package arcsettings

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

type InitializeDisableProcessOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// InitializeDisableProcess ...
func (c ArcSettingsClient) InitializeDisableProcess(ctx context.Context, id ArcSettingId) (result InitializeDisableProcessOperationResponse, err error) {
	req, err := c.preparerForInitializeDisableProcess(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "InitializeDisableProcess", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForInitializeDisableProcess(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "InitializeDisableProcess", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// InitializeDisableProcessThenPoll performs InitializeDisableProcess then polls until it's completed
func (c ArcSettingsClient) InitializeDisableProcessThenPoll(ctx context.Context, id ArcSettingId) error {
	result, err := c.InitializeDisableProcess(ctx, id)
	if err != nil {
		return fmt.Errorf("performing InitializeDisableProcess: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after InitializeDisableProcess: %+v", err)
	}

	return nil
}

// preparerForInitializeDisableProcess prepares the InitializeDisableProcess request.
func (c ArcSettingsClient) preparerForInitializeDisableProcess(ctx context.Context, id ArcSettingId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/initializeDisableProcess", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForInitializeDisableProcess sends the InitializeDisableProcess request. The method will close the
// http.Response Body if it receives an error.
func (c ArcSettingsClient) senderForInitializeDisableProcess(ctx context.Context, req *http.Request) (future InitializeDisableProcessOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
