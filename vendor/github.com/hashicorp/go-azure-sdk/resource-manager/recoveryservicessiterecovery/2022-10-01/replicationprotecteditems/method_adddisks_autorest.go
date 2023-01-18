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

type AddDisksOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// AddDisks ...
func (c ReplicationProtectedItemsClient) AddDisks(ctx context.Context, id ReplicationProtectedItemId, input AddDisksInput) (result AddDisksOperationResponse, err error) {
	req, err := c.preparerForAddDisks(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "AddDisks", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForAddDisks(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "AddDisks", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// AddDisksThenPoll performs AddDisks then polls until it's completed
func (c ReplicationProtectedItemsClient) AddDisksThenPoll(ctx context.Context, id ReplicationProtectedItemId, input AddDisksInput) error {
	result, err := c.AddDisks(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing AddDisks: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after AddDisks: %+v", err)
	}

	return nil
}

// preparerForAddDisks prepares the AddDisks request.
func (c ReplicationProtectedItemsClient) preparerForAddDisks(ctx context.Context, id ReplicationProtectedItemId, input AddDisksInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/addDisks", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForAddDisks sends the AddDisks request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationProtectedItemsClient) senderForAddDisks(ctx context.Context, req *http.Request) (future AddDisksOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
