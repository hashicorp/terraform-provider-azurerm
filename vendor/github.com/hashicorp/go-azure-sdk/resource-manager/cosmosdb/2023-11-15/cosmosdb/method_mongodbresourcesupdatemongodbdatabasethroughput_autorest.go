package cosmosdb

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDBResourcesUpdateMongoDBDatabaseThroughputOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// MongoDBResourcesUpdateMongoDBDatabaseThroughput ...
func (c CosmosDBClient) MongoDBResourcesUpdateMongoDBDatabaseThroughput(ctx context.Context, id MongodbDatabaseId, input ThroughputSettingsUpdateParameters) (result MongoDBResourcesUpdateMongoDBDatabaseThroughputOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesUpdateMongoDBDatabaseThroughput(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesUpdateMongoDBDatabaseThroughput", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMongoDBResourcesUpdateMongoDBDatabaseThroughput(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesUpdateMongoDBDatabaseThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MongoDBResourcesUpdateMongoDBDatabaseThroughputThenPoll performs MongoDBResourcesUpdateMongoDBDatabaseThroughput then polls until it's completed
func (c CosmosDBClient) MongoDBResourcesUpdateMongoDBDatabaseThroughputThenPoll(ctx context.Context, id MongodbDatabaseId, input ThroughputSettingsUpdateParameters) error {
	result, err := c.MongoDBResourcesUpdateMongoDBDatabaseThroughput(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing MongoDBResourcesUpdateMongoDBDatabaseThroughput: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after MongoDBResourcesUpdateMongoDBDatabaseThroughput: %+v", err)
	}

	return nil
}

// preparerForMongoDBResourcesUpdateMongoDBDatabaseThroughput prepares the MongoDBResourcesUpdateMongoDBDatabaseThroughput request.
func (c CosmosDBClient) preparerForMongoDBResourcesUpdateMongoDBDatabaseThroughput(ctx context.Context, id MongodbDatabaseId, input ThroughputSettingsUpdateParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/throughputSettings/default", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForMongoDBResourcesUpdateMongoDBDatabaseThroughput sends the MongoDBResourcesUpdateMongoDBDatabaseThroughput request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForMongoDBResourcesUpdateMongoDBDatabaseThroughput(ctx context.Context, req *http.Request) (future MongoDBResourcesUpdateMongoDBDatabaseThroughputOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
