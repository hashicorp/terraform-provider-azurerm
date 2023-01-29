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

type TestFailoverOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// TestFailover ...
func (c ReplicationProtectedItemsClient) TestFailover(ctx context.Context, id ReplicationProtectedItemId, input TestFailoverInput) (result TestFailoverOperationResponse, err error) {
	req, err := c.preparerForTestFailover(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "TestFailover", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForTestFailover(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "TestFailover", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// TestFailoverThenPoll performs TestFailover then polls until it's completed
func (c ReplicationProtectedItemsClient) TestFailoverThenPoll(ctx context.Context, id ReplicationProtectedItemId, input TestFailoverInput) error {
	result, err := c.TestFailover(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing TestFailover: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after TestFailover: %+v", err)
	}

	return nil
}

// preparerForTestFailover prepares the TestFailover request.
func (c ReplicationProtectedItemsClient) preparerForTestFailover(ctx context.Context, id ReplicationProtectedItemId, input TestFailoverInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/testFailover", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForTestFailover sends the TestFailover request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationProtectedItemsClient) senderForTestFailover(ctx context.Context, req *http.Request) (future TestFailoverOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
