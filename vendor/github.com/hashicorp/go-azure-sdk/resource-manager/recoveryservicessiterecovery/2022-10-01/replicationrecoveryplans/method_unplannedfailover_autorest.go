package replicationrecoveryplans

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

type UnplannedFailoverOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// UnplannedFailover ...
func (c ReplicationRecoveryPlansClient) UnplannedFailover(ctx context.Context, id ReplicationRecoveryPlanId, input RecoveryPlanUnplannedFailoverInput) (result UnplannedFailoverOperationResponse, err error) {
	req, err := c.preparerForUnplannedFailover(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationrecoveryplans.ReplicationRecoveryPlansClient", "UnplannedFailover", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForUnplannedFailover(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationrecoveryplans.ReplicationRecoveryPlansClient", "UnplannedFailover", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// UnplannedFailoverThenPoll performs UnplannedFailover then polls until it's completed
func (c ReplicationRecoveryPlansClient) UnplannedFailoverThenPoll(ctx context.Context, id ReplicationRecoveryPlanId, input RecoveryPlanUnplannedFailoverInput) error {
	result, err := c.UnplannedFailover(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing UnplannedFailover: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after UnplannedFailover: %+v", err)
	}

	return nil
}

// preparerForUnplannedFailover prepares the UnplannedFailover request.
func (c ReplicationRecoveryPlansClient) preparerForUnplannedFailover(ctx context.Context, id ReplicationRecoveryPlanId, input RecoveryPlanUnplannedFailoverInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/unplannedFailover", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForUnplannedFailover sends the UnplannedFailover request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationRecoveryPlansClient) senderForUnplannedFailover(ctx context.Context, req *http.Request) (future UnplannedFailoverOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
