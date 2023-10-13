package accountmigrations

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

type StorageAccountsCustomerInitiatedMigrationOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// StorageAccountsCustomerInitiatedMigration ...
func (c AccountMigrationsClient) StorageAccountsCustomerInitiatedMigration(ctx context.Context, id commonids.StorageAccountId, input StorageAccountMigration) (result StorageAccountsCustomerInitiatedMigrationOperationResponse, err error) {
	req, err := c.preparerForStorageAccountsCustomerInitiatedMigration(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accountmigrations.AccountMigrationsClient", "StorageAccountsCustomerInitiatedMigration", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForStorageAccountsCustomerInitiatedMigration(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accountmigrations.AccountMigrationsClient", "StorageAccountsCustomerInitiatedMigration", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// StorageAccountsCustomerInitiatedMigrationThenPoll performs StorageAccountsCustomerInitiatedMigration then polls until it's completed
func (c AccountMigrationsClient) StorageAccountsCustomerInitiatedMigrationThenPoll(ctx context.Context, id commonids.StorageAccountId, input StorageAccountMigration) error {
	result, err := c.StorageAccountsCustomerInitiatedMigration(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing StorageAccountsCustomerInitiatedMigration: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after StorageAccountsCustomerInitiatedMigration: %+v", err)
	}

	return nil
}

// preparerForStorageAccountsCustomerInitiatedMigration prepares the StorageAccountsCustomerInitiatedMigration request.
func (c AccountMigrationsClient) preparerForStorageAccountsCustomerInitiatedMigration(ctx context.Context, id commonids.StorageAccountId, input StorageAccountMigration) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/startAccountMigration", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForStorageAccountsCustomerInitiatedMigration sends the StorageAccountsCustomerInitiatedMigration request. The method will close the
// http.Response Body if it receives an error.
func (c AccountMigrationsClient) senderForStorageAccountsCustomerInitiatedMigration(ctx context.Context, req *http.Request) (future StorageAccountsCustomerInitiatedMigrationOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
