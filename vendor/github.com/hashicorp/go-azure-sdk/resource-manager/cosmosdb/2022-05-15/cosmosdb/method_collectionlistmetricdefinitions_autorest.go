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

type CollectionListMetricDefinitionsOperationResponse struct {
	HttpResponse *http.Response
	Model        *MetricDefinitionsListResult
}

// CollectionListMetricDefinitions ...
func (c CosmosDBClient) CollectionListMetricDefinitions(ctx context.Context, id CollectionId) (result CollectionListMetricDefinitionsOperationResponse, err error) {
	req, err := c.preparerForCollectionListMetricDefinitions(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionListMetricDefinitions", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionListMetricDefinitions", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCollectionListMetricDefinitions(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionListMetricDefinitions", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCollectionListMetricDefinitions prepares the CollectionListMetricDefinitions request.
func (c CosmosDBClient) preparerForCollectionListMetricDefinitions(ctx context.Context, id CollectionId) (*http.Request, error) {
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

// responderForCollectionListMetricDefinitions handles the response to the CollectionListMetricDefinitions request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForCollectionListMetricDefinitions(resp *http.Response) (result CollectionListMetricDefinitionsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
