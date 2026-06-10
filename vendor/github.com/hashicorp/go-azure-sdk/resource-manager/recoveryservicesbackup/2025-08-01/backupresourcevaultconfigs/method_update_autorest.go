package backupresourcevaultconfigs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *BackupResourceVaultConfigResource
}

type UpdateOperationOptions struct {
	XMsAuthorizationAuxiliary *string
}

func DefaultUpdateOperationOptions() UpdateOperationOptions {
	return UpdateOperationOptions{}
}

func (o UpdateOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.XMsAuthorizationAuxiliary != nil {
		out["x-ms-authorization-auxiliary"] = *o.XMsAuthorizationAuxiliary
	}

	return out
}

func (o UpdateOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// Update ...
func (c BackupResourceVaultConfigsClient) Update(ctx context.Context, id VaultId, input BackupResourceVaultConfigResource, options UpdateOperationOptions) (result UpdateOperationResponse, err error) {
	req, err := c.preparerForUpdate(ctx, id, input, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupresourcevaultconfigs.BackupResourceVaultConfigsClient", "Update", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupresourcevaultconfigs.BackupResourceVaultConfigsClient", "Update", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupresourcevaultconfigs.BackupResourceVaultConfigsClient", "Update", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUpdate prepares the Update request.
func (c BackupResourceVaultConfigsClient) preparerForUpdate(ctx context.Context, id VaultId, input BackupResourceVaultConfigResource, options UpdateOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/backupconfig/vaultconfig", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForUpdate handles the response to the Update request. The method always
// closes the http.Response Body.
func (c BackupResourceVaultConfigsClient) responderForUpdate(resp *http.Response) (result UpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
