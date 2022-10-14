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

type ReprotectOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Reprotect ...
func (c ReplicationProtectedItemsClient) Reprotect(ctx context.Context, id ReplicationProtectedItemId, input ReverseReplicationInput) (result ReprotectOperationResponse, err error) {
	req, err := c.preparerForReprotect(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "Reprotect", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForReprotect(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "Reprotect", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ReprotectThenPoll performs Reprotect then polls until it's completed
func (c ReplicationProtectedItemsClient) ReprotectThenPoll(ctx context.Context, id ReplicationProtectedItemId, input ReverseReplicationInput) error {
	result, err := c.Reprotect(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Reprotect: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Reprotect: %+v", err)
	}

	return nil
}

// preparerForReprotect prepares the Reprotect request.
func (c ReplicationProtectedItemsClient) preparerForReprotect(ctx context.Context, id ReplicationProtectedItemId, input ReverseReplicationInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/reProtect", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForReprotect sends the Reprotect request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationProtectedItemsClient) senderForReprotect(ctx context.Context, req *http.Request) (future ReprotectOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
