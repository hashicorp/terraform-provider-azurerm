package capacities

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

type SuspendOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Suspend ...
func (c CapacitiesClient) Suspend(ctx context.Context, id CapacitiesId) (result SuspendOperationResponse, err error) {
	req, err := c.preparerForSuspend(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "capacities.CapacitiesClient", "Suspend", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSuspend(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "capacities.CapacitiesClient", "Suspend", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SuspendThenPoll performs Suspend then polls until it's completed
func (c CapacitiesClient) SuspendThenPoll(ctx context.Context, id CapacitiesId) error {
	result, err := c.Suspend(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Suspend: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Suspend: %+v", err)
	}

	return nil
}

// preparerForSuspend prepares the Suspend request.
func (c CapacitiesClient) preparerForSuspend(ctx context.Context, id CapacitiesId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/suspend", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForSuspend sends the Suspend request. The method will close the
// http.Response Body if it receives an error.
func (c CapacitiesClient) senderForSuspend(ctx context.Context, req *http.Request) (future SuspendOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
