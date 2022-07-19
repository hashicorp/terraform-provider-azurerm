package redisenterprise

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabasesListKeysOperationResponse struct {
	HttpResponse *http.Response
	Model        *AccessKeys
}

// DatabasesListKeys ...
func (c RedisEnterpriseClient) DatabasesListKeys(ctx context.Context, id DatabaseId) (result DatabasesListKeysOperationResponse, err error) {
	req, err := c.preparerForDatabasesListKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesListKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesListKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDatabasesListKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesListKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDatabasesListKeys prepares the DatabasesListKeys request.
func (c RedisEnterpriseClient) preparerForDatabasesListKeys(ctx context.Context, id DatabaseId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDatabasesListKeys handles the response to the DatabasesListKeys request. The method always
// closes the http.Response Body.
func (c RedisEnterpriseClient) responderForDatabasesListKeys(resp *http.Response) (result DatabasesListKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
