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

type RepairReplicationOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// RepairReplication ...
func (c ReplicationProtectedItemsClient) RepairReplication(ctx context.Context, id ReplicationProtectedItemId) (result RepairReplicationOperationResponse, err error) {
	req, err := c.preparerForRepairReplication(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "RepairReplication", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForRepairReplication(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "RepairReplication", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// RepairReplicationThenPoll performs RepairReplication then polls until it's completed
func (c ReplicationProtectedItemsClient) RepairReplicationThenPoll(ctx context.Context, id ReplicationProtectedItemId) error {
	result, err := c.RepairReplication(ctx, id)
	if err != nil {
		return fmt.Errorf("performing RepairReplication: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after RepairReplication: %+v", err)
	}

	return nil
}

// preparerForRepairReplication prepares the RepairReplication request.
func (c ReplicationProtectedItemsClient) preparerForRepairReplication(ctx context.Context, id ReplicationProtectedItemId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/repairReplication", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForRepairReplication sends the RepairReplication request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationProtectedItemsClient) senderForRepairReplication(ctx context.Context, req *http.Request) (future RepairReplicationOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
