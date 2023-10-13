package encryptionscopes

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PutOperationResponse struct {
	HttpResponse *http.Response
	Model        *EncryptionScope
}

// Put ...
func (c EncryptionScopesClient) Put(ctx context.Context, id EncryptionScopeId, input EncryptionScope) (result PutOperationResponse, err error) {
	req, err := c.preparerForPut(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encryptionscopes.EncryptionScopesClient", "Put", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "encryptionscopes.EncryptionScopesClient", "Put", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPut(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encryptionscopes.EncryptionScopesClient", "Put", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPut prepares the Put request.
func (c EncryptionScopesClient) preparerForPut(ctx context.Context, id EncryptionScopeId, input EncryptionScope) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForPut handles the response to the Put request. The method always
// closes the http.Response Body.
func (c EncryptionScopesClient) responderForPut(resp *http.Response) (result PutOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
