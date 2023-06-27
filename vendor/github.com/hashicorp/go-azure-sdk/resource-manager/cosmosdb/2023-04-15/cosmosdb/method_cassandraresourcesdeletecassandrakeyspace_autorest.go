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

type CassandraResourcesDeleteCassandraKeyspaceOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CassandraResourcesDeleteCassandraKeyspace ...
func (c CosmosDBClient) CassandraResourcesDeleteCassandraKeyspace(ctx context.Context, id CassandraKeyspaceId) (result CassandraResourcesDeleteCassandraKeyspaceOperationResponse, err error) {
	req, err := c.preparerForCassandraResourcesDeleteCassandraKeyspace(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesDeleteCassandraKeyspace", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCassandraResourcesDeleteCassandraKeyspace(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesDeleteCassandraKeyspace", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CassandraResourcesDeleteCassandraKeyspaceThenPoll performs CassandraResourcesDeleteCassandraKeyspace then polls until it's completed
func (c CosmosDBClient) CassandraResourcesDeleteCassandraKeyspaceThenPoll(ctx context.Context, id CassandraKeyspaceId) error {
	result, err := c.CassandraResourcesDeleteCassandraKeyspace(ctx, id)
	if err != nil {
		return fmt.Errorf("performing CassandraResourcesDeleteCassandraKeyspace: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CassandraResourcesDeleteCassandraKeyspace: %+v", err)
	}

	return nil
}

// preparerForCassandraResourcesDeleteCassandraKeyspace prepares the CassandraResourcesDeleteCassandraKeyspace request.
func (c CosmosDBClient) preparerForCassandraResourcesDeleteCassandraKeyspace(ctx context.Context, id CassandraKeyspaceId) (*http.Request, error) {
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

// senderForCassandraResourcesDeleteCassandraKeyspace sends the CassandraResourcesDeleteCassandraKeyspace request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForCassandraResourcesDeleteCassandraKeyspace(ctx context.Context, req *http.Request) (future CassandraResourcesDeleteCassandraKeyspaceOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
