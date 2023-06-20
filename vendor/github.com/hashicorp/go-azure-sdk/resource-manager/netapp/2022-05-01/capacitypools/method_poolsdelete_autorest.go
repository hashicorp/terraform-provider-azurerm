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

type PoolsDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// PoolsDelete ...
func (c CapacityPoolsClient) PoolsDelete(ctx context.Context, id CapacityPoolId) (result PoolsDeleteOperationResponse, err error) {
	req, err := c.preparerForPoolsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "capacitypools.CapacityPoolsClient", "PoolsDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPoolsDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "capacitypools.CapacityPoolsClient", "PoolsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PoolsDeleteThenPoll performs PoolsDelete then polls until it's completed
func (c CapacityPoolsClient) PoolsDeleteThenPoll(ctx context.Context, id CapacityPoolId) error {
	result, err := c.PoolsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing PoolsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after PoolsDelete: %+v", err)
	}

	return nil
}

// preparerForPoolsDelete prepares the PoolsDelete request.
func (c CapacityPoolsClient) preparerForPoolsDelete(ctx context.Context, id CapacityPoolId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForPoolsDelete sends the PoolsDelete request. The method will close the
// http.Response Body if it receives an error.
func (c CapacityPoolsClient) senderForPoolsDelete(ctx context.Context, req *http.Request) (future PoolsDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
