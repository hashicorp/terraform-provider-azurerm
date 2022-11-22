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

type UpdateMobilityServiceOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// UpdateMobilityService ...
func (c ReplicationProtectedItemsClient) UpdateMobilityService(ctx context.Context, id ReplicationProtectedItemId, input UpdateMobilityServiceRequest) (result UpdateMobilityServiceOperationResponse, err error) {
	req, err := c.preparerForUpdateMobilityService(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "UpdateMobilityService", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForUpdateMobilityService(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotecteditems.ReplicationProtectedItemsClient", "UpdateMobilityService", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// UpdateMobilityServiceThenPoll performs UpdateMobilityService then polls until it's completed
func (c ReplicationProtectedItemsClient) UpdateMobilityServiceThenPoll(ctx context.Context, id ReplicationProtectedItemId, input UpdateMobilityServiceRequest) error {
	result, err := c.UpdateMobilityService(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing UpdateMobilityService: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after UpdateMobilityService: %+v", err)
	}

	return nil
}

// preparerForUpdateMobilityService prepares the UpdateMobilityService request.
func (c ReplicationProtectedItemsClient) preparerForUpdateMobilityService(ctx context.Context, id ReplicationProtectedItemId, input UpdateMobilityServiceRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/updateMobilityService", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForUpdateMobilityService sends the UpdateMobilityService request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationProtectedItemsClient) senderForUpdateMobilityService(ctx context.Context, req *http.Request) (future UpdateMobilityServiceOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
