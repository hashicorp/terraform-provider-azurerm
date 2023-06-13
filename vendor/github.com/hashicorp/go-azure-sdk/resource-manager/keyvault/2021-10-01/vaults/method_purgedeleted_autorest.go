package vaults

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

type PurgeDeletedOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// PurgeDeleted ...
func (c VaultsClient) PurgeDeleted(ctx context.Context, id DeletedVaultId) (result PurgeDeletedOperationResponse, err error) {
	req, err := c.preparerForPurgeDeleted(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "PurgeDeleted", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPurgeDeleted(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "PurgeDeleted", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PurgeDeletedThenPoll performs PurgeDeleted then polls until it's completed
func (c VaultsClient) PurgeDeletedThenPoll(ctx context.Context, id DeletedVaultId) error {
	result, err := c.PurgeDeleted(ctx, id)
	if err != nil {
		return fmt.Errorf("performing PurgeDeleted: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after PurgeDeleted: %+v", err)
	}

	return nil
}

// preparerForPurgeDeleted prepares the PurgeDeleted request.
func (c VaultsClient) preparerForPurgeDeleted(ctx context.Context, id DeletedVaultId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/purge", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForPurgeDeleted sends the PurgeDeleted request. The method will close the
// http.Response Body if it receives an error.
func (c VaultsClient) senderForPurgeDeleted(ctx context.Context, req *http.Request) (future PurgeDeletedOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
