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

type CassandraResourcesCreateUpdateCassandraKeyspaceOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CassandraResourcesCreateUpdateCassandraKeyspace ...
func (c CosmosDBClient) CassandraResourcesCreateUpdateCassandraKeyspace(ctx context.Context, id CassandraKeyspaceId, input CassandraKeyspaceCreateUpdateParameters) (result CassandraResourcesCreateUpdateCassandraKeyspaceOperationResponse, err error) {
	req, err := c.preparerForCassandraResourcesCreateUpdateCassandraKeyspace(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesCreateUpdateCassandraKeyspace", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCassandraResourcesCreateUpdateCassandraKeyspace(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesCreateUpdateCassandraKeyspace", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CassandraResourcesCreateUpdateCassandraKeyspaceThenPoll performs CassandraResourcesCreateUpdateCassandraKeyspace then polls until it's completed
func (c CosmosDBClient) CassandraResourcesCreateUpdateCassandraKeyspaceThenPoll(ctx context.Context, id CassandraKeyspaceId, input CassandraKeyspaceCreateUpdateParameters) error {
	result, err := c.CassandraResourcesCreateUpdateCassandraKeyspace(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CassandraResourcesCreateUpdateCassandraKeyspace: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CassandraResourcesCreateUpdateCassandraKeyspace: %+v", err)
	}

	return nil
}

// preparerForCassandraResourcesCreateUpdateCassandraKeyspace prepares the CassandraResourcesCreateUpdateCassandraKeyspace request.
func (c CosmosDBClient) preparerForCassandraResourcesCreateUpdateCassandraKeyspace(ctx context.Context, id CassandraKeyspaceId, input CassandraKeyspaceCreateUpdateParameters) (*http.Request, error) {
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

// senderForCassandraResourcesCreateUpdateCassandraKeyspace sends the CassandraResourcesCreateUpdateCassandraKeyspace request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForCassandraResourcesCreateUpdateCassandraKeyspace(ctx context.Context, req *http.Request) (future CassandraResourcesCreateUpdateCassandraKeyspaceOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
