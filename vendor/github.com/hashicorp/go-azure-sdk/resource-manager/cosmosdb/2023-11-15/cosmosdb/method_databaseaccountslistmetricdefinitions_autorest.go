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

type DatabaseAccountsListMetricDefinitionsOperationResponse struct {
	HttpResponse *http.Response
	Model        *MetricDefinitionsListResult
}

// DatabaseAccountsListMetricDefinitions ...
func (c CosmosDBClient) DatabaseAccountsListMetricDefinitions(ctx context.Context, id DatabaseAccountId) (result DatabaseAccountsListMetricDefinitionsOperationResponse, err error) {
	req, err := c.preparerForDatabaseAccountsListMetricDefinitions(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListMetricDefinitions", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListMetricDefinitions", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDatabaseAccountsListMetricDefinitions(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListMetricDefinitions", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDatabaseAccountsListMetricDefinitions prepares the DatabaseAccountsListMetricDefinitions request.
func (c CosmosDBClient) preparerForDatabaseAccountsListMetricDefinitions(ctx context.Context, id DatabaseAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/metricDefinitions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDatabaseAccountsListMetricDefinitions handles the response to the DatabaseAccountsListMetricDefinitions request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForDatabaseAccountsListMetricDefinitions(resp *http.Response) (result DatabaseAccountsListMetricDefinitionsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
