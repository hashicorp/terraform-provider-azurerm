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

type ListServiceSASOperationResponse struct {
	HttpResponse *http.Response
	Model        *ListServiceSasResponse
}

// ListServiceSAS ...
func (c StorageAccountsClient) ListServiceSAS(ctx context.Context, id commonids.StorageAccountId, input ServiceSasParameters) (result ListServiceSASOperationResponse, err error) {
	req, err := c.preparerForListServiceSAS(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListServiceSAS", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListServiceSAS", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListServiceSAS(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListServiceSAS", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListServiceSAS prepares the ListServiceSAS request.
func (c StorageAccountsClient) preparerForListServiceSAS(ctx context.Context, id commonids.StorageAccountId, input ServiceSasParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listServiceSas", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListServiceSAS handles the response to the ListServiceSAS request. The method always
// closes the http.Response Body.
func (c StorageAccountsClient) responderForListServiceSAS(resp *http.Response) (result ListServiceSASOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
