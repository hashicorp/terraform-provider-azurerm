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

type ListAccountSASOperationResponse struct {
	HttpResponse *http.Response
	Model        *ListAccountSasResponse
}

// ListAccountSAS ...
func (c StorageAccountsClient) ListAccountSAS(ctx context.Context, id commonids.StorageAccountId, input AccountSasParameters) (result ListAccountSASOperationResponse, err error) {
	req, err := c.preparerForListAccountSAS(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListAccountSAS", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListAccountSAS", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListAccountSAS(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListAccountSAS", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListAccountSAS prepares the ListAccountSAS request.
func (c StorageAccountsClient) preparerForListAccountSAS(ctx context.Context, id commonids.StorageAccountId, input AccountSasParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listAccountSas", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListAccountSAS handles the response to the ListAccountSAS request. The method always
// closes the http.Response Body.
func (c StorageAccountsClient) responderForListAccountSAS(resp *http.Response) (result ListAccountSASOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
