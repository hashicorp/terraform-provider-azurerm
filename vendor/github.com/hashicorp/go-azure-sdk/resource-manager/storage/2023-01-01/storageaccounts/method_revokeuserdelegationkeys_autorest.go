package storageaccounts

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

type RevokeUserDelegationKeysOperationResponse struct {
	HttpResponse *http.Response
}

// RevokeUserDelegationKeys ...
func (c StorageAccountsClient) RevokeUserDelegationKeys(ctx context.Context, id commonids.StorageAccountId) (result RevokeUserDelegationKeysOperationResponse, err error) {
	req, err := c.preparerForRevokeUserDelegationKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "RevokeUserDelegationKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "RevokeUserDelegationKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRevokeUserDelegationKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "RevokeUserDelegationKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRevokeUserDelegationKeys prepares the RevokeUserDelegationKeys request.
func (c StorageAccountsClient) preparerForRevokeUserDelegationKeys(ctx context.Context, id commonids.StorageAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/revokeUserDelegationKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRevokeUserDelegationKeys handles the response to the RevokeUserDelegationKeys request. The method always
// closes the http.Response Body.
func (c StorageAccountsClient) responderForRevokeUserDelegationKeys(resp *http.Response) (result RevokeUserDelegationKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
