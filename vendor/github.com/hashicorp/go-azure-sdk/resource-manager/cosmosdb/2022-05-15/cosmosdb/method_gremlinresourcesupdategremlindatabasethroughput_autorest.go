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

type GremlinResourcesUpdateGremlinDatabaseThroughputOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// GremlinResourcesUpdateGremlinDatabaseThroughput ...
func (c CosmosDBClient) GremlinResourcesUpdateGremlinDatabaseThroughput(ctx context.Context, id GremlinDatabaseId, input ThroughputSettingsUpdateParameters) (result GremlinResourcesUpdateGremlinDatabaseThroughputOperationResponse, err error) {
	req, err := c.preparerForGremlinResourcesUpdateGremlinDatabaseThroughput(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesUpdateGremlinDatabaseThroughput", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForGremlinResourcesUpdateGremlinDatabaseThroughput(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesUpdateGremlinDatabaseThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// GremlinResourcesUpdateGremlinDatabaseThroughputThenPoll performs GremlinResourcesUpdateGremlinDatabaseThroughput then polls until it's completed
func (c CosmosDBClient) GremlinResourcesUpdateGremlinDatabaseThroughputThenPoll(ctx context.Context, id GremlinDatabaseId, input ThroughputSettingsUpdateParameters) error {
	result, err := c.GremlinResourcesUpdateGremlinDatabaseThroughput(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing GremlinResourcesUpdateGremlinDatabaseThroughput: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after GremlinResourcesUpdateGremlinDatabaseThroughput: %+v", err)
	}

	return nil
}

// preparerForGremlinResourcesUpdateGremlinDatabaseThroughput prepares the GremlinResourcesUpdateGremlinDatabaseThroughput request.
func (c CosmosDBClient) preparerForGremlinResourcesUpdateGremlinDatabaseThroughput(ctx context.Context, id GremlinDatabaseId, input ThroughputSettingsUpdateParameters) (*http.Request, error) {
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

// senderForGremlinResourcesUpdateGremlinDatabaseThroughput sends the GremlinResourcesUpdateGremlinDatabaseThroughput request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForGremlinResourcesUpdateGremlinDatabaseThroughput(ctx context.Context, req *http.Request) (future GremlinResourcesUpdateGremlinDatabaseThroughputOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
