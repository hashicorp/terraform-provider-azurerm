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

type FailoverCommitOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// FailoverCommit ...
func (c ReplicationProtectedItemsClient) FailoverCommit(ctx context.Context, id ReplicationProtectedItemId) (result FailoverCommitOperationResponse, err error) {
	req, err := c.preparerForFailoverCommit(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "FailoverCommit", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForFailoverCommit(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "FailoverCommit", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// FailoverCommitThenPoll performs FailoverCommit then polls until it's completed
func (c ReplicationProtectedItemsClient) FailoverCommitThenPoll(ctx context.Context, id ReplicationProtectedItemId) error {
	result, err := c.FailoverCommit(ctx, id)
	if err != nil {
		return fmt.Errorf("performing FailoverCommit: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after FailoverCommit: %+v", err)
	}

	return nil
}

// preparerForFailoverCommit prepares the FailoverCommit request.
func (c ReplicationProtectedItemsClient) preparerForFailoverCommit(ctx context.Context, id ReplicationProtectedItemId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/failoverCommit", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForFailoverCommit sends the FailoverCommit request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationProtectedItemsClient) senderForFailoverCommit(ctx context.Context, req *http.Request) (future FailoverCommitOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
