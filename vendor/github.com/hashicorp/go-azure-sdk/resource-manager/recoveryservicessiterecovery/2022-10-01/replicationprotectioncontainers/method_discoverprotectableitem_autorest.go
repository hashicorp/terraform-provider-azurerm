package replicationprotectioncontainers

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

type DiscoverProtectableItemOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DiscoverProtectableItem ...
func (c ReplicationProtectionContainersClient) DiscoverProtectableItem(ctx context.Context, id ReplicationProtectionContainerId, input DiscoverProtectableItemRequest) (result DiscoverProtectableItemOperationResponse, err error) {
	req, err := c.preparerForDiscoverProtectableItem(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotectioncontainers.ReplicationProtectionContainersClient", "DiscoverProtectableItem", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDiscoverProtectableItem(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotectioncontainers.ReplicationProtectionContainersClient", "DiscoverProtectableItem", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DiscoverProtectableItemThenPoll performs DiscoverProtectableItem then polls until it's completed
func (c ReplicationProtectionContainersClient) DiscoverProtectableItemThenPoll(ctx context.Context, id ReplicationProtectionContainerId, input DiscoverProtectableItemRequest) error {
	result, err := c.DiscoverProtectableItem(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DiscoverProtectableItem: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DiscoverProtectableItem: %+v", err)
	}

	return nil
}

// preparerForDiscoverProtectableItem prepares the DiscoverProtectableItem request.
func (c ReplicationProtectionContainersClient) preparerForDiscoverProtectableItem(ctx context.Context, id ReplicationProtectionContainerId, input DiscoverProtectableItemRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/discoverProtectableItem", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDiscoverProtectableItem sends the DiscoverProtectableItem request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationProtectionContainersClient) senderForDiscoverProtectableItem(ctx context.Context, req *http.Request) (future DiscoverProtectableItemOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
