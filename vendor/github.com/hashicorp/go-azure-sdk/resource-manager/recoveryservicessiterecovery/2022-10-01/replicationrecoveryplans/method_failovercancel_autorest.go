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

type FailoverCancelOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// FailoverCancel ...
func (c ReplicationRecoveryPlansClient) FailoverCancel(ctx context.Context, id ReplicationRecoveryPlanId) (result FailoverCancelOperationResponse, err error) {
	req, err := c.preparerForFailoverCancel(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationrecoveryplans.ReplicationRecoveryPlansClient", "FailoverCancel", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForFailoverCancel(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationrecoveryplans.ReplicationRecoveryPlansClient", "FailoverCancel", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// FailoverCancelThenPoll performs FailoverCancel then polls until it's completed
func (c ReplicationRecoveryPlansClient) FailoverCancelThenPoll(ctx context.Context, id ReplicationRecoveryPlanId) error {
	result, err := c.FailoverCancel(ctx, id)
	if err != nil {
		return fmt.Errorf("performing FailoverCancel: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after FailoverCancel: %+v", err)
	}

	return nil
}

// preparerForFailoverCancel prepares the FailoverCancel request.
func (c ReplicationRecoveryPlansClient) preparerForFailoverCancel(ctx context.Context, id ReplicationRecoveryPlanId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/failoverCancel", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForFailoverCancel sends the FailoverCancel request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationRecoveryPlansClient) senderForFailoverCancel(ctx context.Context, req *http.Request) (future FailoverCancelOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
