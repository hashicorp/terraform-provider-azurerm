package lab

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

type SyncGroupOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SyncGroup ...
func (c LabClient) SyncGroup(ctx context.Context, id LabId) (result SyncGroupOperationResponse, err error) {
	req, err := c.preparerForSyncGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "lab.LabClient", "SyncGroup", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSyncGroup(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "lab.LabClient", "SyncGroup", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SyncGroupThenPoll performs SyncGroup then polls until it's completed
func (c LabClient) SyncGroupThenPoll(ctx context.Context, id LabId) error {
	result, err := c.SyncGroup(ctx, id)
	if err != nil {
		return fmt.Errorf("performing SyncGroup: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SyncGroup: %+v", err)
	}

	return nil
}

// preparerForSyncGroup prepares the SyncGroup request.
func (c LabClient) preparerForSyncGroup(ctx context.Context, id LabId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/syncGroup", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForSyncGroup sends the SyncGroup request. The method will close the
// http.Response Body if it receives an error.
func (c LabClient) senderForSyncGroup(ctx context.Context, req *http.Request) (future SyncGroupOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
