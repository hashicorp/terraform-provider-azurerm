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

type MigrateToAadOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// MigrateToAad ...
func (c ReplicationFabricsClient) MigrateToAad(ctx context.Context, id ReplicationFabricId) (result MigrateToAadOperationResponse, err error) {
	req, err := c.preparerForMigrateToAad(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationfabrics.ReplicationFabricsClient", "MigrateToAad", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMigrateToAad(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationfabrics.ReplicationFabricsClient", "MigrateToAad", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MigrateToAadThenPoll performs MigrateToAad then polls until it's completed
func (c ReplicationFabricsClient) MigrateToAadThenPoll(ctx context.Context, id ReplicationFabricId) error {
	result, err := c.MigrateToAad(ctx, id)
	if err != nil {
		return fmt.Errorf("performing MigrateToAad: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after MigrateToAad: %+v", err)
	}

	return nil
}

// preparerForMigrateToAad prepares the MigrateToAad request.
func (c ReplicationFabricsClient) preparerForMigrateToAad(ctx context.Context, id ReplicationFabricId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/migratetoaad", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForMigrateToAad sends the MigrateToAad request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationFabricsClient) senderForMigrateToAad(ctx context.Context, req *http.Request) (future MigrateToAadOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
