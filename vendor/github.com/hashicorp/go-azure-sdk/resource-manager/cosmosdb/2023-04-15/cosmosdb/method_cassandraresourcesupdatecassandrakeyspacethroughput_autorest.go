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

type CassandraResourcesUpdateCassandraKeyspaceThroughputOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CassandraResourcesUpdateCassandraKeyspaceThroughput ...
func (c CosmosDBClient) CassandraResourcesUpdateCassandraKeyspaceThroughput(ctx context.Context, id CassandraKeyspaceId, input ThroughputSettingsUpdateParameters) (result CassandraResourcesUpdateCassandraKeyspaceThroughputOperationResponse, err error) {
	req, err := c.preparerForCassandraResourcesUpdateCassandraKeyspaceThroughput(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesUpdateCassandraKeyspaceThroughput", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCassandraResourcesUpdateCassandraKeyspaceThroughput(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesUpdateCassandraKeyspaceThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CassandraResourcesUpdateCassandraKeyspaceThroughputThenPoll performs CassandraResourcesUpdateCassandraKeyspaceThroughput then polls until it's completed
func (c CosmosDBClient) CassandraResourcesUpdateCassandraKeyspaceThroughputThenPoll(ctx context.Context, id CassandraKeyspaceId, input ThroughputSettingsUpdateParameters) error {
	result, err := c.CassandraResourcesUpdateCassandraKeyspaceThroughput(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CassandraResourcesUpdateCassandraKeyspaceThroughput: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CassandraResourcesUpdateCassandraKeyspaceThroughput: %+v", err)
	}

	return nil
}

// preparerForCassandraResourcesUpdateCassandraKeyspaceThroughput prepares the CassandraResourcesUpdateCassandraKeyspaceThroughput request.
func (c CosmosDBClient) preparerForCassandraResourcesUpdateCassandraKeyspaceThroughput(ctx context.Context, id CassandraKeyspaceId, input ThroughputSettingsUpdateParameters) (*http.Request, error) {
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

// senderForCassandraResourcesUpdateCassandraKeyspaceThroughput sends the CassandraResourcesUpdateCassandraKeyspaceThroughput request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForCassandraResourcesUpdateCassandraKeyspaceThroughput(ctx context.Context, req *http.Request) (future CassandraResourcesUpdateCassandraKeyspaceThroughputOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
