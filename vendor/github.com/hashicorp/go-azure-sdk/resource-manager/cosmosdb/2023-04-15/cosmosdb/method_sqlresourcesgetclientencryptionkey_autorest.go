package cosmosdb

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlResourcesGetClientEncryptionKeyOperationResponse struct {
	HttpResponse *http.Response
	Model        *ClientEncryptionKeyGetResults
}

// SqlResourcesGetClientEncryptionKey ...
func (c CosmosDBClient) SqlResourcesGetClientEncryptionKey(ctx context.Context, id ClientEncryptionKeyId) (result SqlResourcesGetClientEncryptionKeyOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesGetClientEncryptionKey(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesGetClientEncryptionKey", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesGetClientEncryptionKey", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSqlResourcesGetClientEncryptionKey(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesGetClientEncryptionKey", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSqlResourcesGetClientEncryptionKey prepares the SqlResourcesGetClientEncryptionKey request.
func (c CosmosDBClient) preparerForSqlResourcesGetClientEncryptionKey(ctx context.Context, id ClientEncryptionKeyId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSqlResourcesGetClientEncryptionKey handles the response to the SqlResourcesGetClientEncryptionKey request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForSqlResourcesGetClientEncryptionKey(resp *http.Response) (result SqlResourcesGetClientEncryptionKeyOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
