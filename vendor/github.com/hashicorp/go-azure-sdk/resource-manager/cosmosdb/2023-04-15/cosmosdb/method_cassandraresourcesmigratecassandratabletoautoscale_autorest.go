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

type CassandraResourcesMigrateCassandraTableToAutoscaleOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CassandraResourcesMigrateCassandraTableToAutoscale ...
func (c CosmosDBClient) CassandraResourcesMigrateCassandraTableToAutoscale(ctx context.Context, id CassandraKeyspaceTableId) (result CassandraResourcesMigrateCassandraTableToAutoscaleOperationResponse, err error) {
	req, err := c.preparerForCassandraResourcesMigrateCassandraTableToAutoscale(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesMigrateCassandraTableToAutoscale", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCassandraResourcesMigrateCassandraTableToAutoscale(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesMigrateCassandraTableToAutoscale", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CassandraResourcesMigrateCassandraTableToAutoscaleThenPoll performs CassandraResourcesMigrateCassandraTableToAutoscale then polls until it's completed
func (c CosmosDBClient) CassandraResourcesMigrateCassandraTableToAutoscaleThenPoll(ctx context.Context, id CassandraKeyspaceTableId) error {
	result, err := c.CassandraResourcesMigrateCassandraTableToAutoscale(ctx, id)
	if err != nil {
		return fmt.Errorf("performing CassandraResourcesMigrateCassandraTableToAutoscale: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CassandraResourcesMigrateCassandraTableToAutoscale: %+v", err)
	}

	return nil
}

// preparerForCassandraResourcesMigrateCassandraTableToAutoscale prepares the CassandraResourcesMigrateCassandraTableToAutoscale request.
func (c CosmosDBClient) preparerForCassandraResourcesMigrateCassandraTableToAutoscale(ctx context.Context, id CassandraKeyspaceTableId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/throughputSettings/default/migrateToAutoscale", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForCassandraResourcesMigrateCassandraTableToAutoscale sends the CassandraResourcesMigrateCassandraTableToAutoscale request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForCassandraResourcesMigrateCassandraTableToAutoscale(ctx context.Context, req *http.Request) (future CassandraResourcesMigrateCassandraTableToAutoscaleOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
