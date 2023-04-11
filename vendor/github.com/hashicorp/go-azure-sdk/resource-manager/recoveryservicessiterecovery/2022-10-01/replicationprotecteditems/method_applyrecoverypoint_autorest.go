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

type ApplyRecoveryPointOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ApplyRecoveryPoint ...
func (c ReplicationProtectedItemsClient) ApplyRecoveryPoint(ctx context.Context, id ReplicationProtectedItemId, input ApplyRecoveryPointInput) (result ApplyRecoveryPointOperationResponse, err error) {
	req, err := c.preparerForApplyRecoveryPoint(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "ApplyRecoveryPoint", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForApplyRecoveryPoint(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "ApplyRecoveryPoint", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ApplyRecoveryPointThenPoll performs ApplyRecoveryPoint then polls until it's completed
func (c ReplicationProtectedItemsClient) ApplyRecoveryPointThenPoll(ctx context.Context, id ReplicationProtectedItemId, input ApplyRecoveryPointInput) error {
	result, err := c.ApplyRecoveryPoint(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ApplyRecoveryPoint: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ApplyRecoveryPoint: %+v", err)
	}

	return nil
}

// preparerForApplyRecoveryPoint prepares the ApplyRecoveryPoint request.
func (c ReplicationProtectedItemsClient) preparerForApplyRecoveryPoint(ctx context.Context, id ReplicationProtectedItemId, input ApplyRecoveryPointInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/applyRecoveryPoint", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForApplyRecoveryPoint sends the ApplyRecoveryPoint request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationProtectedItemsClient) senderForApplyRecoveryPoint(ctx context.Context, req *http.Request) (future ApplyRecoveryPointOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
