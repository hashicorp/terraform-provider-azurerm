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

type CassandraResourcesCreateUpdateCassandraTableOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CassandraResourcesCreateUpdateCassandraTable ...
func (c CosmosDBClient) CassandraResourcesCreateUpdateCassandraTable(ctx context.Context, id CassandraKeyspaceTableId, input CassandraTableCreateUpdateParameters) (result CassandraResourcesCreateUpdateCassandraTableOperationResponse, err error) {
	req, err := c.preparerForCassandraResourcesCreateUpdateCassandraTable(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesCreateUpdateCassandraTable", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCassandraResourcesCreateUpdateCassandraTable(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesCreateUpdateCassandraTable", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CassandraResourcesCreateUpdateCassandraTableThenPoll performs CassandraResourcesCreateUpdateCassandraTable then polls until it's completed
func (c CosmosDBClient) CassandraResourcesCreateUpdateCassandraTableThenPoll(ctx context.Context, id CassandraKeyspaceTableId, input CassandraTableCreateUpdateParameters) error {
	result, err := c.CassandraResourcesCreateUpdateCassandraTable(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CassandraResourcesCreateUpdateCassandraTable: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CassandraResourcesCreateUpdateCassandraTable: %+v", err)
	}

	return nil
}

// preparerForCassandraResourcesCreateUpdateCassandraTable prepares the CassandraResourcesCreateUpdateCassandraTable request.
func (c CosmosDBClient) preparerForCassandraResourcesCreateUpdateCassandraTable(ctx context.Context, id CassandraKeyspaceTableId, input CassandraTableCreateUpdateParameters) (*http.Request, error) {
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

// senderForCassandraResourcesCreateUpdateCassandraTable sends the CassandraResourcesCreateUpdateCassandraTable request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForCassandraResourcesCreateUpdateCassandraTable(ctx context.Context, req *http.Request) (future CassandraResourcesCreateUpdateCassandraTableOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
