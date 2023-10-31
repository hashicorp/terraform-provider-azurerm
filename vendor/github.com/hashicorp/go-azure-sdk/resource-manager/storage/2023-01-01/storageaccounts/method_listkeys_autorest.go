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

type ListKeysOperationResponse struct {
	HttpResponse *http.Response
	Model        *StorageAccountListKeysResult
}

type ListKeysOperationOptions struct {
	Expand *ListKeyExpand
}

func DefaultListKeysOperationOptions() ListKeysOperationOptions {
	return ListKeysOperationOptions{}
}

func (o ListKeysOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListKeysOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Expand != nil {
		out["$expand"] = *o.Expand
	}

	return out
}

// ListKeys ...
func (c StorageAccountsClient) ListKeys(ctx context.Context, id commonids.StorageAccountId, options ListKeysOperationOptions) (result ListKeysOperationResponse, err error) {
	req, err := c.preparerForListKeys(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListKeys prepares the ListKeys request.
func (c StorageAccountsClient) preparerForListKeys(ctx context.Context, id commonids.StorageAccountId, options ListKeysOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/listKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListKeys handles the response to the ListKeys request. The method always
// closes the http.Response Body.
func (c StorageAccountsClient) responderForListKeys(resp *http.Response) (result ListKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
