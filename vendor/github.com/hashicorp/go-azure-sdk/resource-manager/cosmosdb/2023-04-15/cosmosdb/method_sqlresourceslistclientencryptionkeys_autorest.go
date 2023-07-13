package cosmosdb

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlResourcesListClientEncryptionKeysOperationResponse struct {
	HttpResponse *http.Response
	Model        *ClientEncryptionKeysListResult
}

// SqlResourcesListClientEncryptionKeys ...
func (c CosmosDBClient) SqlResourcesListClientEncryptionKeys(ctx context.Context, id SqlDatabaseId) (result SqlResourcesListClientEncryptionKeysOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesListClientEncryptionKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesListClientEncryptionKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesListClientEncryptionKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSqlResourcesListClientEncryptionKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesListClientEncryptionKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSqlResourcesListClientEncryptionKeys prepares the SqlResourcesListClientEncryptionKeys request.
func (c CosmosDBClient) preparerForSqlResourcesListClientEncryptionKeys(ctx context.Context, id SqlDatabaseId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/clientEncryptionKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSqlResourcesListClientEncryptionKeys handles the response to the SqlResourcesListClientEncryptionKeys request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForSqlResourcesListClientEncryptionKeys(resp *http.Response) (result SqlResourcesListClientEncryptionKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
