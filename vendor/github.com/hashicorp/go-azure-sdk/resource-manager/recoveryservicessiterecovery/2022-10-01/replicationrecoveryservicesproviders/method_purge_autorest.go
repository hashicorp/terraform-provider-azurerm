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

type PurgeOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Purge ...
func (c ReplicationRecoveryServicesProvidersClient) Purge(ctx context.Context, id ReplicationRecoveryServicesProviderId) (result PurgeOperationResponse, err error) {
	req, err := c.preparerForPurge(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationrecoveryservicesproviders.ReplicationRecoveryServicesProvidersClient", "Purge", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPurge(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationrecoveryservicesproviders.ReplicationRecoveryServicesProvidersClient", "Purge", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PurgeThenPoll performs Purge then polls until it's completed
func (c ReplicationRecoveryServicesProvidersClient) PurgeThenPoll(ctx context.Context, id ReplicationRecoveryServicesProviderId) error {
	result, err := c.Purge(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Purge: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Purge: %+v", err)
	}

	return nil
}

// preparerForPurge prepares the Purge request.
func (c ReplicationRecoveryServicesProvidersClient) preparerForPurge(ctx context.Context, id ReplicationRecoveryServicesProviderId) (*http.Request, error) {
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

// senderForPurge sends the Purge request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationRecoveryServicesProvidersClient) senderForPurge(ctx context.Context, req *http.Request) (future PurgeOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
