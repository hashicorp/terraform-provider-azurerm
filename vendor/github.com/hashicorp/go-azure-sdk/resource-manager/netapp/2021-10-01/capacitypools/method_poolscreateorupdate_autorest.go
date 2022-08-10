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

type PoolsCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// PoolsCreateOrUpdate ...
func (c CapacityPoolsClient) PoolsCreateOrUpdate(ctx context.Context, id CapacityPoolId, input CapacityPool) (result PoolsCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForPoolsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "capacitypools.CapacityPoolsClient", "PoolsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPoolsCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "capacitypools.CapacityPoolsClient", "PoolsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PoolsCreateOrUpdateThenPoll performs PoolsCreateOrUpdate then polls until it's completed
func (c CapacityPoolsClient) PoolsCreateOrUpdateThenPoll(ctx context.Context, id CapacityPoolId, input CapacityPool) error {
	result, err := c.PoolsCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing PoolsCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after PoolsCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForPoolsCreateOrUpdate prepares the PoolsCreateOrUpdate request.
func (c CapacityPoolsClient) preparerForPoolsCreateOrUpdate(ctx context.Context, id CapacityPoolId, input CapacityPool) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForPoolsCreateOrUpdate sends the PoolsCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c CapacityPoolsClient) senderForPoolsCreateOrUpdate(ctx context.Context, req *http.Request) (future PoolsCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
