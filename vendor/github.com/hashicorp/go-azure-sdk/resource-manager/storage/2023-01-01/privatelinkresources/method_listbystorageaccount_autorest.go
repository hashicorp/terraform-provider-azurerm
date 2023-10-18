package privatelinkresources

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

type ListByStorageAccountOperationResponse struct {
	HttpResponse *http.Response
	Model        *PrivateLinkResourceListResult
}

// ListByStorageAccount ...
func (c PrivateLinkResourcesClient) ListByStorageAccount(ctx context.Context, id commonids.StorageAccountId) (result ListByStorageAccountOperationResponse, err error) {
	req, err := c.preparerForListByStorageAccount(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkresources.PrivateLinkResourcesClient", "ListByStorageAccount", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkresources.PrivateLinkResourcesClient", "ListByStorageAccount", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListByStorageAccount(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkresources.PrivateLinkResourcesClient", "ListByStorageAccount", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListByStorageAccount prepares the ListByStorageAccount request.
func (c PrivateLinkResourcesClient) preparerForListByStorageAccount(ctx context.Context, id commonids.StorageAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/privateLinkResources", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListByStorageAccount handles the response to the ListByStorageAccount request. The method always
// closes the http.Response Body.
func (c PrivateLinkResourcesClient) responderForListByStorageAccount(resp *http.Response) (result ListByStorageAccountOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
