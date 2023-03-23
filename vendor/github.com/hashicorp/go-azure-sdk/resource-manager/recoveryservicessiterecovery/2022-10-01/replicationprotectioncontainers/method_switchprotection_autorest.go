package replicationprotectioncontainers

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

type SwitchProtectionOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SwitchProtection ...
func (c ReplicationProtectionContainersClient) SwitchProtection(ctx context.Context, id ReplicationProtectionContainerId, input SwitchProtectionInput) (result SwitchProtectionOperationResponse, err error) {
	req, err := c.preparerForSwitchProtection(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotectioncontainers.ReplicationProtectionContainersClient", "SwitchProtection", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSwitchProtection(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotectioncontainers.ReplicationProtectionContainersClient", "SwitchProtection", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SwitchProtectionThenPoll performs SwitchProtection then polls until it's completed
func (c ReplicationProtectionContainersClient) SwitchProtectionThenPoll(ctx context.Context, id ReplicationProtectionContainerId, input SwitchProtectionInput) error {
	result, err := c.SwitchProtection(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing SwitchProtection: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SwitchProtection: %+v", err)
	}

	return nil
}

// preparerForSwitchProtection prepares the SwitchProtection request.
func (c ReplicationProtectionContainersClient) preparerForSwitchProtection(ctx context.Context, id ReplicationProtectionContainerId, input SwitchProtectionInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/switchprotection", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForSwitchProtection sends the SwitchProtection request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationProtectionContainersClient) senderForSwitchProtection(ctx context.Context, req *http.Request) (future SwitchProtectionOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
