package storageaccounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AbortHierarchicalNamespaceMigrationOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// AbortHierarchicalNamespaceMigration ...
func (c StorageAccountsClient) AbortHierarchicalNamespaceMigration(ctx context.Context, id commonids.StorageAccountId) (result AbortHierarchicalNamespaceMigrationOperationResponse, err error) {
	req, err := c.preparerForAbortHierarchicalNamespaceMigration(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "AbortHierarchicalNamespaceMigration", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForAbortHierarchicalNamespaceMigration(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "AbortHierarchicalNamespaceMigration", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// AbortHierarchicalNamespaceMigrationThenPoll performs AbortHierarchicalNamespaceMigration then polls until it's completed
func (c StorageAccountsClient) AbortHierarchicalNamespaceMigrationThenPoll(ctx context.Context, id commonids.StorageAccountId) error {
	result, err := c.AbortHierarchicalNamespaceMigration(ctx, id)
	if err != nil {
		return fmt.Errorf("performing AbortHierarchicalNamespaceMigration: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after AbortHierarchicalNamespaceMigration: %+v", err)
	}

	return nil
}

// preparerForAbortHierarchicalNamespaceMigration prepares the AbortHierarchicalNamespaceMigration request.
func (c StorageAccountsClient) preparerForAbortHierarchicalNamespaceMigration(ctx context.Context, id commonids.StorageAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/aborthnsonmigration", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForAbortHierarchicalNamespaceMigration sends the AbortHierarchicalNamespaceMigration request. The method will close the
// http.Response Body if it receives an error.
func (c StorageAccountsClient) senderForAbortHierarchicalNamespaceMigration(ctx context.Context, req *http.Request) (future AbortHierarchicalNamespaceMigrationOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
