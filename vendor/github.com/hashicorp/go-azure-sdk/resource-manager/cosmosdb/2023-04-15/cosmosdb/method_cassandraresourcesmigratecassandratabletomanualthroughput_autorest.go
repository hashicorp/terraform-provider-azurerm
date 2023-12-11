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

type CassandraResourcesMigrateCassandraTableToManualThroughputOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CassandraResourcesMigrateCassandraTableToManualThroughput ...
func (c CosmosDBClient) CassandraResourcesMigrateCassandraTableToManualThroughput(ctx context.Context, id CassandraKeyspaceTableId) (result CassandraResourcesMigrateCassandraTableToManualThroughputOperationResponse, err error) {
	req, err := c.preparerForCassandraResourcesMigrateCassandraTableToManualThroughput(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesMigrateCassandraTableToManualThroughput", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCassandraResourcesMigrateCassandraTableToManualThroughput(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesMigrateCassandraTableToManualThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CassandraResourcesMigrateCassandraTableToManualThroughputThenPoll performs CassandraResourcesMigrateCassandraTableToManualThroughput then polls until it's completed
func (c CosmosDBClient) CassandraResourcesMigrateCassandraTableToManualThroughputThenPoll(ctx context.Context, id CassandraKeyspaceTableId) error {
	result, err := c.CassandraResourcesMigrateCassandraTableToManualThroughput(ctx, id)
	if err != nil {
		return fmt.Errorf("performing CassandraResourcesMigrateCassandraTableToManualThroughput: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CassandraResourcesMigrateCassandraTableToManualThroughput: %+v", err)
	}

	return nil
}

// preparerForCassandraResourcesMigrateCassandraTableToManualThroughput prepares the CassandraResourcesMigrateCassandraTableToManualThroughput request.
func (c CosmosDBClient) preparerForCassandraResourcesMigrateCassandraTableToManualThroughput(ctx context.Context, id CassandraKeyspaceTableId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/throughputSettings/default/migrateToManualThroughput", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForCassandraResourcesMigrateCassandraTableToManualThroughput sends the CassandraResourcesMigrateCassandraTableToManualThroughput request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForCassandraResourcesMigrateCassandraTableToManualThroughput(ctx context.Context, req *http.Request) (future CassandraResourcesMigrateCassandraTableToManualThroughputOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
