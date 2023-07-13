package caches

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

type SpaceAllocationOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SpaceAllocation ...
func (c CachesClient) SpaceAllocation(ctx context.Context, id CacheId, input []StorageTargetSpaceAllocation) (result SpaceAllocationOperationResponse, err error) {
	req, err := c.preparerForSpaceAllocation(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "caches.CachesClient", "SpaceAllocation", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSpaceAllocation(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "caches.CachesClient", "SpaceAllocation", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SpaceAllocationThenPoll performs SpaceAllocation then polls until it's completed
func (c CachesClient) SpaceAllocationThenPoll(ctx context.Context, id CacheId, input []StorageTargetSpaceAllocation) error {
	result, err := c.SpaceAllocation(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing SpaceAllocation: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SpaceAllocation: %+v", err)
	}

	return nil
}

// preparerForSpaceAllocation prepares the SpaceAllocation request.
func (c CachesClient) preparerForSpaceAllocation(ctx context.Context, id CacheId, input []StorageTargetSpaceAllocation) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/spaceAllocation", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForSpaceAllocation sends the SpaceAllocation request. The method will close the
// http.Response Body if it receives an error.
func (c CachesClient) senderForSpaceAllocation(ctx context.Context, req *http.Request) (future SpaceAllocationOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
