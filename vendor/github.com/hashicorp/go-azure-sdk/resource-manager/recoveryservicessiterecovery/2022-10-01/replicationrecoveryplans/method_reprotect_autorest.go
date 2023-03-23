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

type ReprotectOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Reprotect ...
func (c ReplicationRecoveryPlansClient) Reprotect(ctx context.Context, id ReplicationRecoveryPlanId) (result ReprotectOperationResponse, err error) {
	req, err := c.preparerForReprotect(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationrecoveryplans.ReplicationRecoveryPlansClient", "Reprotect", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForReprotect(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationrecoveryplans.ReplicationRecoveryPlansClient", "Reprotect", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ReprotectThenPoll performs Reprotect then polls until it's completed
func (c ReplicationRecoveryPlansClient) ReprotectThenPoll(ctx context.Context, id ReplicationRecoveryPlanId) error {
	result, err := c.Reprotect(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Reprotect: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Reprotect: %+v", err)
	}

	return nil
}

// preparerForReprotect prepares the Reprotect request.
func (c ReplicationRecoveryPlansClient) preparerForReprotect(ctx context.Context, id ReplicationRecoveryPlanId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/reProtect", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForReprotect sends the Reprotect request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationRecoveryPlansClient) senderForReprotect(ctx context.Context, req *http.Request) (future ReprotectOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
