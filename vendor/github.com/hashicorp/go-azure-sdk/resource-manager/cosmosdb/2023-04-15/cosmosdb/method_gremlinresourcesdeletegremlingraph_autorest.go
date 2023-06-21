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

type GremlinResourcesDeleteGremlinGraphOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// GremlinResourcesDeleteGremlinGraph ...
func (c CosmosDBClient) GremlinResourcesDeleteGremlinGraph(ctx context.Context, id GraphId) (result GremlinResourcesDeleteGremlinGraphOperationResponse, err error) {
	req, err := c.preparerForGremlinResourcesDeleteGremlinGraph(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesDeleteGremlinGraph", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForGremlinResourcesDeleteGremlinGraph(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesDeleteGremlinGraph", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// GremlinResourcesDeleteGremlinGraphThenPoll performs GremlinResourcesDeleteGremlinGraph then polls until it's completed
func (c CosmosDBClient) GremlinResourcesDeleteGremlinGraphThenPoll(ctx context.Context, id GraphId) error {
	result, err := c.GremlinResourcesDeleteGremlinGraph(ctx, id)
	if err != nil {
		return fmt.Errorf("performing GremlinResourcesDeleteGremlinGraph: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after GremlinResourcesDeleteGremlinGraph: %+v", err)
	}

	return nil
}

// preparerForGremlinResourcesDeleteGremlinGraph prepares the GremlinResourcesDeleteGremlinGraph request.
func (c CosmosDBClient) preparerForGremlinResourcesDeleteGremlinGraph(ctx context.Context, id GraphId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForGremlinResourcesDeleteGremlinGraph sends the GremlinResourcesDeleteGremlinGraph request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForGremlinResourcesDeleteGremlinGraph(ctx context.Context, req *http.Request) (future GremlinResourcesDeleteGremlinGraphOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
