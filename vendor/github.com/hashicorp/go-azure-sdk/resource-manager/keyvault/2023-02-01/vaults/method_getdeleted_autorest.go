package vaults

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetDeletedOperationResponse struct {
	HttpResponse *http.Response
	Model        *DeletedVault
}

// GetDeleted ...
func (c VaultsClient) GetDeleted(ctx context.Context, id DeletedVaultId) (result GetDeletedOperationResponse, err error) {
	req, err := c.preparerForGetDeleted(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "GetDeleted", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "GetDeleted", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetDeleted(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "GetDeleted", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetDeleted prepares the GetDeleted request.
func (c VaultsClient) preparerForGetDeleted(ctx context.Context, id DeletedVaultId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetDeleted handles the response to the GetDeleted request. The method always
// closes the http.Response Body.
func (c VaultsClient) responderForGetDeleted(resp *http.Response) (result GetDeletedOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
