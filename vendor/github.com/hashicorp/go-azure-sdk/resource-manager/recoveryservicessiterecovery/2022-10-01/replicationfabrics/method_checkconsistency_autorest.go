package replicationfabrics

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

type CheckConsistencyOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CheckConsistency ...
func (c ReplicationFabricsClient) CheckConsistency(ctx context.Context, id ReplicationFabricId) (result CheckConsistencyOperationResponse, err error) {
	req, err := c.preparerForCheckConsistency(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationfabrics.ReplicationFabricsClient", "CheckConsistency", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCheckConsistency(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationfabrics.ReplicationFabricsClient", "CheckConsistency", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CheckConsistencyThenPoll performs CheckConsistency then polls until it's completed
func (c ReplicationFabricsClient) CheckConsistencyThenPoll(ctx context.Context, id ReplicationFabricId) error {
	result, err := c.CheckConsistency(ctx, id)
	if err != nil {
		return fmt.Errorf("performing CheckConsistency: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CheckConsistency: %+v", err)
	}

	return nil
}

// preparerForCheckConsistency prepares the CheckConsistency request.
func (c ReplicationFabricsClient) preparerForCheckConsistency(ctx context.Context, id ReplicationFabricId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/checkConsistency", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForCheckConsistency sends the CheckConsistency request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationFabricsClient) senderForCheckConsistency(ctx context.Context, req *http.Request) (future CheckConsistencyOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
