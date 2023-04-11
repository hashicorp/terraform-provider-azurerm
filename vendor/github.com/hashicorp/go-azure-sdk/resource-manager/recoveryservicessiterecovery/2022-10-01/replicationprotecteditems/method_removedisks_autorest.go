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

type RemoveDisksOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// RemoveDisks ...
func (c ReplicationProtectedItemsClient) RemoveDisks(ctx context.Context, id ReplicationProtectedItemId, input RemoveDisksInput) (result RemoveDisksOperationResponse, err error) {
	req, err := c.preparerForRemoveDisks(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "RemoveDisks", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForRemoveDisks(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "RemoveDisks", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// RemoveDisksThenPoll performs RemoveDisks then polls until it's completed
func (c ReplicationProtectedItemsClient) RemoveDisksThenPoll(ctx context.Context, id ReplicationProtectedItemId, input RemoveDisksInput) error {
	result, err := c.RemoveDisks(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing RemoveDisks: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after RemoveDisks: %+v", err)
	}

	return nil
}

// preparerForRemoveDisks prepares the RemoveDisks request.
func (c ReplicationProtectedItemsClient) preparerForRemoveDisks(ctx context.Context, id ReplicationProtectedItemId, input RemoveDisksInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/removeDisks", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForRemoveDisks sends the RemoveDisks request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationProtectedItemsClient) senderForRemoveDisks(ctx context.Context, req *http.Request) (future RemoveDisksOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
