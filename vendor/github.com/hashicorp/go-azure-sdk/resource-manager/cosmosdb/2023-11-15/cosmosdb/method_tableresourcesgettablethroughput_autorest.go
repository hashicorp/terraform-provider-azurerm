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

type TableResourcesGetTableThroughputOperationResponse struct {
	HttpResponse *http.Response
	Model        *ThroughputSettingsGetResults
}

// TableResourcesGetTableThroughput ...
func (c CosmosDBClient) TableResourcesGetTableThroughput(ctx context.Context, id TableId) (result TableResourcesGetTableThroughputOperationResponse, err error) {
	req, err := c.preparerForTableResourcesGetTableThroughput(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "TableResourcesGetTableThroughput", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "TableResourcesGetTableThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTableResourcesGetTableThroughput(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "TableResourcesGetTableThroughput", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTableResourcesGetTableThroughput prepares the TableResourcesGetTableThroughput request.
func (c CosmosDBClient) preparerForTableResourcesGetTableThroughput(ctx context.Context, id TableId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/throughputSettings/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForTableResourcesGetTableThroughput handles the response to the TableResourcesGetTableThroughput request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForTableResourcesGetTableThroughput(resp *http.Response) (result TableResourcesGetTableThroughputOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
