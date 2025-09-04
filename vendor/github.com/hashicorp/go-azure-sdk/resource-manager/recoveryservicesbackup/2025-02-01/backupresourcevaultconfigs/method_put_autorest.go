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

type PutOperationResponse struct {
	HttpResponse *http.Response
	Model        *BackupResourceVaultConfigResource
}

type PutOperationOptions struct {
	XMsAuthorizationAuxiliary *string
}

func DefaultPutOperationOptions() PutOperationOptions {
	return PutOperationOptions{}
}

func (o PutOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.XMsAuthorizationAuxiliary != nil {
		out["x-ms-authorization-auxiliary"] = *o.XMsAuthorizationAuxiliary
	}

	return out
}

func (o PutOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// Put ...
func (c BackupResourceVaultConfigsClient) Put(ctx context.Context, id VaultId, input BackupResourceVaultConfigResource, options PutOperationOptions) (result PutOperationResponse, err error) {
	req, err := c.preparerForPut(ctx, id, input, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupresourcevaultconfigs.BackupResourceVaultConfigsClient", "Put", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupresourcevaultconfigs.BackupResourceVaultConfigsClient", "Put", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPut(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupresourcevaultconfigs.BackupResourceVaultConfigsClient", "Put", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPut prepares the Put request.
func (c BackupResourceVaultConfigsClient) preparerForPut(ctx context.Context, id VaultId, input BackupResourceVaultConfigResource, options PutOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/backupconfig/vaultconfig", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForPut handles the response to the Put request. The method always
// closes the http.Response Body.
func (c BackupResourceVaultConfigsClient) responderForPut(resp *http.Response) (result PutOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
