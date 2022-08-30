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

type GremlinResourcesUpdateGremlinGraphThroughputOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// GremlinResourcesUpdateGremlinGraphThroughput ...
func (c CosmosDBClient) GremlinResourcesUpdateGremlinGraphThroughput(ctx context.Context, id GraphId, input ThroughputSettingsUpdateParameters) (result GremlinResourcesUpdateGremlinGraphThroughputOperationResponse, err error) {
	req, err := c.preparerForGremlinResourcesUpdateGremlinGraphThroughput(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesUpdateGremlinGraphThroughput", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForGremlinResourcesUpdateGremlinGraphThroughput(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesUpdateGremlinGraphThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// GremlinResourcesUpdateGremlinGraphThroughputThenPoll performs GremlinResourcesUpdateGremlinGraphThroughput then polls until it's completed
func (c CosmosDBClient) GremlinResourcesUpdateGremlinGraphThroughputThenPoll(ctx context.Context, id GraphId, input ThroughputSettingsUpdateParameters) error {
	result, err := c.GremlinResourcesUpdateGremlinGraphThroughput(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing GremlinResourcesUpdateGremlinGraphThroughput: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after GremlinResourcesUpdateGremlinGraphThroughput: %+v", err)
	}

	return nil
}

// preparerForGremlinResourcesUpdateGremlinGraphThroughput prepares the GremlinResourcesUpdateGremlinGraphThroughput request.
func (c CosmosDBClient) preparerForGremlinResourcesUpdateGremlinGraphThroughput(ctx context.Context, id GraphId, input ThroughputSettingsUpdateParameters) (*http.Request, error) {
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

// senderForGremlinResourcesUpdateGremlinGraphThroughput sends the GremlinResourcesUpdateGremlinGraphThroughput request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForGremlinResourcesUpdateGremlinGraphThroughput(ctx context.Context, req *http.Request) (future GremlinResourcesUpdateGremlinGraphThroughputOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
