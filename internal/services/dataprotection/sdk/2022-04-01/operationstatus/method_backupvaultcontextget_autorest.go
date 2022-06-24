package operationstatus

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupVaultContextGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *OperationResource
}

// BackupVaultContextGet ...
func (c OperationStatusClient) BackupVaultContextGet(ctx context.Context, id BackupVaultOperationStatuId) (result BackupVaultContextGetOperationResponse, err error) {
	req, err := c.preparerForBackupVaultContextGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationstatus.OperationStatusClient", "BackupVaultContextGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationstatus.OperationStatusClient", "BackupVaultContextGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForBackupVaultContextGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationstatus.OperationStatusClient", "BackupVaultContextGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForBackupVaultContextGet prepares the BackupVaultContextGet request.
func (c OperationStatusClient) preparerForBackupVaultContextGet(ctx context.Context, id BackupVaultOperationStatuId) (*http.Request, error) {
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

// responderForBackupVaultContextGet handles the response to the BackupVaultContextGet request. The method always
// closes the http.Response Body.
func (c OperationStatusClient) responderForBackupVaultContextGet(resp *http.Response) (result BackupVaultContextGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
