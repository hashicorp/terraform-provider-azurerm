package capacitypools

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

type PoolsUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// PoolsUpdate ...
func (c CapacityPoolsClient) PoolsUpdate(ctx context.Context, id CapacityPoolId, input CapacityPoolPatch) (result PoolsUpdateOperationResponse, err error) {
	req, err := c.preparerForPoolsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "capacitypools.CapacityPoolsClient", "PoolsUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPoolsUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "capacitypools.CapacityPoolsClient", "PoolsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PoolsUpdateThenPoll performs PoolsUpdate then polls until it's completed
func (c CapacityPoolsClient) PoolsUpdateThenPoll(ctx context.Context, id CapacityPoolId, input CapacityPoolPatch) error {
	result, err := c.PoolsUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing PoolsUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after PoolsUpdate: %+v", err)
	}

	return nil
}

// preparerForPoolsUpdate prepares the PoolsUpdate request.
func (c CapacityPoolsClient) preparerForPoolsUpdate(ctx context.Context, id CapacityPoolId, input CapacityPoolPatch) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForPoolsUpdate sends the PoolsUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c CapacityPoolsClient) senderForPoolsUpdate(ctx context.Context, req *http.Request) (future PoolsUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
