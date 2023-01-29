package accounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MediaservicesSyncStorageKeysOperationResponse struct {
	HttpResponse *http.Response
}

// MediaservicesSyncStorageKeys ...
func (c AccountsClient) MediaservicesSyncStorageKeys(ctx context.Context, id MediaServiceId, input SyncStorageKeysInput) (result MediaservicesSyncStorageKeysOperationResponse, err error) {
	req, err := c.preparerForMediaservicesSyncStorageKeys(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesSyncStorageKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesSyncStorageKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMediaservicesSyncStorageKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesSyncStorageKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMediaservicesSyncStorageKeys prepares the MediaservicesSyncStorageKeys request.
func (c AccountsClient) preparerForMediaservicesSyncStorageKeys(ctx context.Context, id MediaServiceId, input SyncStorageKeysInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/syncStorageKeys", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForMediaservicesSyncStorageKeys handles the response to the MediaservicesSyncStorageKeys request. The method always
// closes the http.Response Body.
func (c AccountsClient) responderForMediaservicesSyncStorageKeys(resp *http.Response) (result MediaservicesSyncStorageKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
