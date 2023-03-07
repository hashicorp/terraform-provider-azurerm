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

type MongoDBResourcesUpdateMongoDBCollectionThroughputOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// MongoDBResourcesUpdateMongoDBCollectionThroughput ...
func (c CosmosDBClient) MongoDBResourcesUpdateMongoDBCollectionThroughput(ctx context.Context, id MongodbDatabaseCollectionId, input ThroughputSettingsUpdateParameters) (result MongoDBResourcesUpdateMongoDBCollectionThroughputOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesUpdateMongoDBCollectionThroughput(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesUpdateMongoDBCollectionThroughput", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMongoDBResourcesUpdateMongoDBCollectionThroughput(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesUpdateMongoDBCollectionThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MongoDBResourcesUpdateMongoDBCollectionThroughputThenPoll performs MongoDBResourcesUpdateMongoDBCollectionThroughput then polls until it's completed
func (c CosmosDBClient) MongoDBResourcesUpdateMongoDBCollectionThroughputThenPoll(ctx context.Context, id MongodbDatabaseCollectionId, input ThroughputSettingsUpdateParameters) error {
	result, err := c.MongoDBResourcesUpdateMongoDBCollectionThroughput(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing MongoDBResourcesUpdateMongoDBCollectionThroughput: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after MongoDBResourcesUpdateMongoDBCollectionThroughput: %+v", err)
	}

	return nil
}

// preparerForMongoDBResourcesUpdateMongoDBCollectionThroughput prepares the MongoDBResourcesUpdateMongoDBCollectionThroughput request.
func (c CosmosDBClient) preparerForMongoDBResourcesUpdateMongoDBCollectionThroughput(ctx context.Context, id MongodbDatabaseCollectionId, input ThroughputSettingsUpdateParameters) (*http.Request, error) {
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

// senderForMongoDBResourcesUpdateMongoDBCollectionThroughput sends the MongoDBResourcesUpdateMongoDBCollectionThroughput request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForMongoDBResourcesUpdateMongoDBCollectionThroughput(ctx context.Context, req *http.Request) (future MongoDBResourcesUpdateMongoDBCollectionThroughputOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
