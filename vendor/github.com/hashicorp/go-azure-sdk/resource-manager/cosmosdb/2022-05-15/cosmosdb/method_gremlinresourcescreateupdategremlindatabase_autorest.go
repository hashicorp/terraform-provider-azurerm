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

type GremlinResourcesCreateUpdateGremlinDatabaseOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// GremlinResourcesCreateUpdateGremlinDatabase ...
func (c CosmosDBClient) GremlinResourcesCreateUpdateGremlinDatabase(ctx context.Context, id GremlinDatabaseId, input GremlinDatabaseCreateUpdateParameters) (result GremlinResourcesCreateUpdateGremlinDatabaseOperationResponse, err error) {
	req, err := c.preparerForGremlinResourcesCreateUpdateGremlinDatabase(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesCreateUpdateGremlinDatabase", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForGremlinResourcesCreateUpdateGremlinDatabase(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesCreateUpdateGremlinDatabase", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// GremlinResourcesCreateUpdateGremlinDatabaseThenPoll performs GremlinResourcesCreateUpdateGremlinDatabase then polls until it's completed
func (c CosmosDBClient) GremlinResourcesCreateUpdateGremlinDatabaseThenPoll(ctx context.Context, id GremlinDatabaseId, input GremlinDatabaseCreateUpdateParameters) error {
	result, err := c.GremlinResourcesCreateUpdateGremlinDatabase(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing GremlinResourcesCreateUpdateGremlinDatabase: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after GremlinResourcesCreateUpdateGremlinDatabase: %+v", err)
	}

	return nil
}

// preparerForGremlinResourcesCreateUpdateGremlinDatabase prepares the GremlinResourcesCreateUpdateGremlinDatabase request.
func (c CosmosDBClient) preparerForGremlinResourcesCreateUpdateGremlinDatabase(ctx context.Context, id GremlinDatabaseId, input GremlinDatabaseCreateUpdateParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForGremlinResourcesCreateUpdateGremlinDatabase sends the GremlinResourcesCreateUpdateGremlinDatabase request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForGremlinResourcesCreateUpdateGremlinDatabase(ctx context.Context, req *http.Request) (future GremlinResourcesCreateUpdateGremlinDatabaseOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
