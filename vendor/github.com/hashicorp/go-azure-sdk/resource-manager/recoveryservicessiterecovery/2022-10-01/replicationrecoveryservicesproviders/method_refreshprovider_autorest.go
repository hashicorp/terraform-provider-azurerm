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

type RefreshProviderOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// RefreshProvider ...
func (c ReplicationRecoveryServicesProvidersClient) RefreshProvider(ctx context.Context, id ReplicationRecoveryServicesProviderId) (result RefreshProviderOperationResponse, err error) {
	req, err := c.preparerForRefreshProvider(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationrecoveryservicesproviders.ReplicationRecoveryServicesProvidersClient", "RefreshProvider", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForRefreshProvider(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationrecoveryservicesproviders.ReplicationRecoveryServicesProvidersClient", "RefreshProvider", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// RefreshProviderThenPoll performs RefreshProvider then polls until it's completed
func (c ReplicationRecoveryServicesProvidersClient) RefreshProviderThenPoll(ctx context.Context, id ReplicationRecoveryServicesProviderId) error {
	result, err := c.RefreshProvider(ctx, id)
	if err != nil {
		return fmt.Errorf("performing RefreshProvider: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after RefreshProvider: %+v", err)
	}

	return nil
}

// preparerForRefreshProvider prepares the RefreshProvider request.
func (c ReplicationRecoveryServicesProvidersClient) preparerForRefreshProvider(ctx context.Context, id ReplicationRecoveryServicesProviderId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/refreshProvider", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForRefreshProvider sends the RefreshProvider request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationRecoveryServicesProvidersClient) senderForRefreshProvider(ctx context.Context, req *http.Request) (future RefreshProviderOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
