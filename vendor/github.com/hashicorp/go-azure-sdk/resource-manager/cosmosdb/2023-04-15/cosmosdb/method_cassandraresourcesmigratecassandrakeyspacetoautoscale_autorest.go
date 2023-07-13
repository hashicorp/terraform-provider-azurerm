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

type CassandraResourcesMigrateCassandraKeyspaceToAutoscaleOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CassandraResourcesMigrateCassandraKeyspaceToAutoscale ...
func (c CosmosDBClient) CassandraResourcesMigrateCassandraKeyspaceToAutoscale(ctx context.Context, id CassandraKeyspaceId) (result CassandraResourcesMigrateCassandraKeyspaceToAutoscaleOperationResponse, err error) {
	req, err := c.preparerForCassandraResourcesMigrateCassandraKeyspaceToAutoscale(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesMigrateCassandraKeyspaceToAutoscale", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCassandraResourcesMigrateCassandraKeyspaceToAutoscale(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesMigrateCassandraKeyspaceToAutoscale", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CassandraResourcesMigrateCassandraKeyspaceToAutoscaleThenPoll performs CassandraResourcesMigrateCassandraKeyspaceToAutoscale then polls until it's completed
func (c CosmosDBClient) CassandraResourcesMigrateCassandraKeyspaceToAutoscaleThenPoll(ctx context.Context, id CassandraKeyspaceId) error {
	result, err := c.CassandraResourcesMigrateCassandraKeyspaceToAutoscale(ctx, id)
	if err != nil {
		return fmt.Errorf("performing CassandraResourcesMigrateCassandraKeyspaceToAutoscale: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CassandraResourcesMigrateCassandraKeyspaceToAutoscale: %+v", err)
	}

	return nil
}

// preparerForCassandraResourcesMigrateCassandraKeyspaceToAutoscale prepares the CassandraResourcesMigrateCassandraKeyspaceToAutoscale request.
func (c CosmosDBClient) preparerForCassandraResourcesMigrateCassandraKeyspaceToAutoscale(ctx context.Context, id CassandraKeyspaceId) (*http.Request, error) {
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

// senderForCassandraResourcesMigrateCassandraKeyspaceToAutoscale sends the CassandraResourcesMigrateCassandraKeyspaceToAutoscale request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForCassandraResourcesMigrateCassandraKeyspaceToAutoscale(ctx context.Context, req *http.Request) (future CassandraResourcesMigrateCassandraKeyspaceToAutoscaleOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
