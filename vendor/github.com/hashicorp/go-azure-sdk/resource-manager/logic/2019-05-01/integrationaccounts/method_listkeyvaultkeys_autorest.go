package integrationaccounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListKeyVaultKeysOperationResponse struct {
	HttpResponse *http.Response
	Model        *KeyVaultKeyCollection
}

// ListKeyVaultKeys ...
func (c IntegrationAccountsClient) ListKeyVaultKeys(ctx context.Context, id IntegrationAccountId, input ListKeyVaultKeysDefinition) (result ListKeyVaultKeysOperationResponse, err error) {
	req, err := c.preparerForListKeyVaultKeys(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "integrationaccounts.IntegrationAccountsClient", "ListKeyVaultKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "integrationaccounts.IntegrationAccountsClient", "ListKeyVaultKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListKeyVaultKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "integrationaccounts.IntegrationAccountsClient", "ListKeyVaultKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListKeyVaultKeys prepares the ListKeyVaultKeys request.
func (c IntegrationAccountsClient) preparerForListKeyVaultKeys(ctx context.Context, id IntegrationAccountId, input ListKeyVaultKeysDefinition) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listKeyVaultKeys", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListKeyVaultKeys handles the response to the ListKeyVaultKeys request. The method always
// closes the http.Response Body.
func (c IntegrationAccountsClient) responderForListKeyVaultKeys(resp *http.Response) (result ListKeyVaultKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
