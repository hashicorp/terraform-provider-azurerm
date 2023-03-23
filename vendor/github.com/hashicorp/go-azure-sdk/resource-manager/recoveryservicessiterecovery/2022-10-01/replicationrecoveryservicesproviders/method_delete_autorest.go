package replicationrecoveryservicesproviders

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

type DeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Delete ...
func (c ReplicationRecoveryServicesProvidersClient) Delete(ctx context.Context, id ReplicationRecoveryServicesProviderId) (result DeleteOperationResponse, err error) {
	req, err := c.preparerForDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationrecoveryservicesproviders.ReplicationRecoveryServicesProvidersClient", "Delete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationrecoveryservicesproviders.ReplicationRecoveryServicesProvidersClient", "Delete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DeleteThenPoll performs Delete then polls until it's completed
func (c ReplicationRecoveryServicesProvidersClient) DeleteThenPoll(ctx context.Context, id ReplicationRecoveryServicesProviderId) error {
	result, err := c.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Delete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Delete: %+v", err)
	}

	return nil
}

// preparerForDelete prepares the Delete request.
func (c ReplicationRecoveryServicesProvidersClient) preparerForDelete(ctx context.Context, id ReplicationRecoveryServicesProviderId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/remove", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDelete sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationRecoveryServicesProvidersClient) senderForDelete(ctx context.Context, req *http.Request) (future DeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
