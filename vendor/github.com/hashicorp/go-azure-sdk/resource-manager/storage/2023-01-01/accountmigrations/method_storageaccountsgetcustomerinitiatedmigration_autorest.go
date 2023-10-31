package accountmigrations

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageAccountsGetCustomerInitiatedMigrationOperationResponse struct {
	HttpResponse *http.Response
	Model        *StorageAccountMigration
}

// StorageAccountsGetCustomerInitiatedMigration ...
func (c AccountMigrationsClient) StorageAccountsGetCustomerInitiatedMigration(ctx context.Context, id commonids.StorageAccountId) (result StorageAccountsGetCustomerInitiatedMigrationOperationResponse, err error) {
	req, err := c.preparerForStorageAccountsGetCustomerInitiatedMigration(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accountmigrations.AccountMigrationsClient", "StorageAccountsGetCustomerInitiatedMigration", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accountmigrations.AccountMigrationsClient", "StorageAccountsGetCustomerInitiatedMigration", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForStorageAccountsGetCustomerInitiatedMigration(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accountmigrations.AccountMigrationsClient", "StorageAccountsGetCustomerInitiatedMigration", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForStorageAccountsGetCustomerInitiatedMigration prepares the StorageAccountsGetCustomerInitiatedMigration request.
func (c AccountMigrationsClient) preparerForStorageAccountsGetCustomerInitiatedMigration(ctx context.Context, id commonids.StorageAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/accountMigrations/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForStorageAccountsGetCustomerInitiatedMigration handles the response to the StorageAccountsGetCustomerInitiatedMigration request. The method always
// closes the http.Response Body.
func (c AccountMigrationsClient) responderForStorageAccountsGetCustomerInitiatedMigration(resp *http.Response) (result StorageAccountsGetCustomerInitiatedMigrationOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
