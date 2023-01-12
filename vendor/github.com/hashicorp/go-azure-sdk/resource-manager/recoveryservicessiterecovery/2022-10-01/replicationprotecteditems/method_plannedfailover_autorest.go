package replicationprotecteditems

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

type PlannedFailoverOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// PlannedFailover ...
func (c ReplicationProtectedItemsClient) PlannedFailover(ctx context.Context, id ReplicationProtectedItemId, input PlannedFailoverInput) (result PlannedFailoverOperationResponse, err error) {
	req, err := c.preparerForPlannedFailover(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "PlannedFailover", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPlannedFailover(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "PlannedFailover", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PlannedFailoverThenPoll performs PlannedFailover then polls until it's completed
func (c ReplicationProtectedItemsClient) PlannedFailoverThenPoll(ctx context.Context, id ReplicationProtectedItemId, input PlannedFailoverInput) error {
	result, err := c.PlannedFailover(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing PlannedFailover: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after PlannedFailover: %+v", err)
	}

	return nil
}

// preparerForPlannedFailover prepares the PlannedFailover request.
func (c ReplicationProtectedItemsClient) preparerForPlannedFailover(ctx context.Context, id ReplicationProtectedItemId, input PlannedFailoverInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/plannedFailover", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForPlannedFailover sends the PlannedFailover request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationProtectedItemsClient) senderForPlannedFailover(ctx context.Context, req *http.Request) (future PlannedFailoverOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
