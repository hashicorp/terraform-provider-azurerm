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

type CassandraResourcesDeleteCassandraTableOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CassandraResourcesDeleteCassandraTable ...
func (c CosmosDBClient) CassandraResourcesDeleteCassandraTable(ctx context.Context, id CassandraKeyspaceTableId) (result CassandraResourcesDeleteCassandraTableOperationResponse, err error) {
	req, err := c.preparerForCassandraResourcesDeleteCassandraTable(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesDeleteCassandraTable", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCassandraResourcesDeleteCassandraTable(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesDeleteCassandraTable", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CassandraResourcesDeleteCassandraTableThenPoll performs CassandraResourcesDeleteCassandraTable then polls until it's completed
func (c CosmosDBClient) CassandraResourcesDeleteCassandraTableThenPoll(ctx context.Context, id CassandraKeyspaceTableId) error {
	result, err := c.CassandraResourcesDeleteCassandraTable(ctx, id)
	if err != nil {
		return fmt.Errorf("performing CassandraResourcesDeleteCassandraTable: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CassandraResourcesDeleteCassandraTable: %+v", err)
	}

	return nil
}

// preparerForCassandraResourcesDeleteCassandraTable prepares the CassandraResourcesDeleteCassandraTable request.
func (c CosmosDBClient) preparerForCassandraResourcesDeleteCassandraTable(ctx context.Context, id CassandraKeyspaceTableId) (*http.Request, error) {
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

// senderForCassandraResourcesDeleteCassandraTable sends the CassandraResourcesDeleteCassandraTable request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForCassandraResourcesDeleteCassandraTable(ctx context.Context, req *http.Request) (future CassandraResourcesDeleteCassandraTableOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
