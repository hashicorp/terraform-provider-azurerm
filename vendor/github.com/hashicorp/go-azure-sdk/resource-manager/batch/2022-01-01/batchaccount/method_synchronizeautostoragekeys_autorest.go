package batchaccount

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SynchronizeAutoStorageKeysOperationResponse struct {
	HttpResponse *http.Response
}

// SynchronizeAutoStorageKeys ...
func (c BatchAccountClient) SynchronizeAutoStorageKeys(ctx context.Context, id BatchAccountId) (result SynchronizeAutoStorageKeysOperationResponse, err error) {
	req, err := c.preparerForSynchronizeAutoStorageKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "batchaccount.BatchAccountClient", "SynchronizeAutoStorageKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "batchaccount.BatchAccountClient", "SynchronizeAutoStorageKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSynchronizeAutoStorageKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "batchaccount.BatchAccountClient", "SynchronizeAutoStorageKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSynchronizeAutoStorageKeys prepares the SynchronizeAutoStorageKeys request.
func (c BatchAccountClient) preparerForSynchronizeAutoStorageKeys(ctx context.Context, id BatchAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/syncAutoStorageKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSynchronizeAutoStorageKeys handles the response to the SynchronizeAutoStorageKeys request. The method always
// closes the http.Response Body.
func (c BatchAccountClient) responderForSynchronizeAutoStorageKeys(resp *http.Response) (result SynchronizeAutoStorageKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
