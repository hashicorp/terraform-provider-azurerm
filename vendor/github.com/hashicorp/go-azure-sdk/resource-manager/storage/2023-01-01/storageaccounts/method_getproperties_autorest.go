package storageaccounts

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetPropertiesOperationResponse struct {
	HttpResponse *http.Response
	Model        *StorageAccount
}

type GetPropertiesOperationOptions struct {
	Expand *StorageAccountExpand
}

func DefaultGetPropertiesOperationOptions() GetPropertiesOperationOptions {
	return GetPropertiesOperationOptions{}
}

func (o GetPropertiesOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o GetPropertiesOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Expand != nil {
		out["$expand"] = *o.Expand
	}

	return out
}

// GetProperties ...
func (c StorageAccountsClient) GetProperties(ctx context.Context, id commonids.StorageAccountId, options GetPropertiesOperationOptions) (result GetPropertiesOperationResponse, err error) {
	req, err := c.preparerForGetProperties(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "GetProperties", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "GetProperties", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetProperties(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "GetProperties", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetProperties prepares the GetProperties request.
func (c StorageAccountsClient) preparerForGetProperties(ctx context.Context, id commonids.StorageAccountId, options GetPropertiesOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetProperties handles the response to the GetProperties request. The method always
// closes the http.Response Body.
func (c StorageAccountsClient) responderForGetProperties(resp *http.Response) (result GetPropertiesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
