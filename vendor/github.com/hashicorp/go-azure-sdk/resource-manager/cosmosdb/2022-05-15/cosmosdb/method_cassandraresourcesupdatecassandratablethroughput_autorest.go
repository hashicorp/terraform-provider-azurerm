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

type CassandraResourcesUpdateCassandraTableThroughputOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CassandraResourcesUpdateCassandraTableThroughput ...
func (c CosmosDBClient) CassandraResourcesUpdateCassandraTableThroughput(ctx context.Context, id CassandraKeyspaceTableId, input ThroughputSettingsUpdateParameters) (result CassandraResourcesUpdateCassandraTableThroughputOperationResponse, err error) {
	req, err := c.preparerForCassandraResourcesUpdateCassandraTableThroughput(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesUpdateCassandraTableThroughput", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCassandraResourcesUpdateCassandraTableThroughput(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesUpdateCassandraTableThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CassandraResourcesUpdateCassandraTableThroughputThenPoll performs CassandraResourcesUpdateCassandraTableThroughput then polls until it's completed
func (c CosmosDBClient) CassandraResourcesUpdateCassandraTableThroughputThenPoll(ctx context.Context, id CassandraKeyspaceTableId, input ThroughputSettingsUpdateParameters) error {
	result, err := c.CassandraResourcesUpdateCassandraTableThroughput(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CassandraResourcesUpdateCassandraTableThroughput: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CassandraResourcesUpdateCassandraTableThroughput: %+v", err)
	}

	return nil
}

// preparerForCassandraResourcesUpdateCassandraTableThroughput prepares the CassandraResourcesUpdateCassandraTableThroughput request.
func (c CosmosDBClient) preparerForCassandraResourcesUpdateCassandraTableThroughput(ctx context.Context, id CassandraKeyspaceTableId, input ThroughputSettingsUpdateParameters) (*http.Request, error) {
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

// senderForCassandraResourcesUpdateCassandraTableThroughput sends the CassandraResourcesUpdateCassandraTableThroughput request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForCassandraResourcesUpdateCassandraTableThroughput(ctx context.Context, req *http.Request) (future CassandraResourcesUpdateCassandraTableThroughputOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
