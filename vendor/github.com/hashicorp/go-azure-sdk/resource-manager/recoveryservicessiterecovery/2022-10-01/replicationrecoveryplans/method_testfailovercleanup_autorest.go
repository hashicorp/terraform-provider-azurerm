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

type TestFailoverCleanupOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// TestFailoverCleanup ...
func (c ReplicationRecoveryPlansClient) TestFailoverCleanup(ctx context.Context, id ReplicationRecoveryPlanId, input RecoveryPlanTestFailoverCleanupInput) (result TestFailoverCleanupOperationResponse, err error) {
	req, err := c.preparerForTestFailoverCleanup(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationrecoveryplans.ReplicationRecoveryPlansClient", "TestFailoverCleanup", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForTestFailoverCleanup(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationrecoveryplans.ReplicationRecoveryPlansClient", "TestFailoverCleanup", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// TestFailoverCleanupThenPoll performs TestFailoverCleanup then polls until it's completed
func (c ReplicationRecoveryPlansClient) TestFailoverCleanupThenPoll(ctx context.Context, id ReplicationRecoveryPlanId, input RecoveryPlanTestFailoverCleanupInput) error {
	result, err := c.TestFailoverCleanup(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing TestFailoverCleanup: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after TestFailoverCleanup: %+v", err)
	}

	return nil
}

// preparerForTestFailoverCleanup prepares the TestFailoverCleanup request.
func (c ReplicationRecoveryPlansClient) preparerForTestFailoverCleanup(ctx context.Context, id ReplicationRecoveryPlanId, input RecoveryPlanTestFailoverCleanupInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/testFailoverCleanup", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForTestFailoverCleanup sends the TestFailoverCleanup request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationRecoveryPlansClient) senderForTestFailoverCleanup(ctx context.Context, req *http.Request) (future TestFailoverCleanupOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
