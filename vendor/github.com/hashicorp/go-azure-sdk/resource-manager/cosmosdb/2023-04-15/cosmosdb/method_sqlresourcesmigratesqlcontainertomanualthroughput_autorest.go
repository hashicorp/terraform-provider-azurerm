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

type SqlResourcesMigrateSqlContainerToManualThroughputOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SqlResourcesMigrateSqlContainerToManualThroughput ...
func (c CosmosDBClient) SqlResourcesMigrateSqlContainerToManualThroughput(ctx context.Context, id ContainerId) (result SqlResourcesMigrateSqlContainerToManualThroughputOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesMigrateSqlContainerToManualThroughput(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesMigrateSqlContainerToManualThroughput", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSqlResourcesMigrateSqlContainerToManualThroughput(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesMigrateSqlContainerToManualThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SqlResourcesMigrateSqlContainerToManualThroughputThenPoll performs SqlResourcesMigrateSqlContainerToManualThroughput then polls until it's completed
func (c CosmosDBClient) SqlResourcesMigrateSqlContainerToManualThroughputThenPoll(ctx context.Context, id ContainerId) error {
	result, err := c.SqlResourcesMigrateSqlContainerToManualThroughput(ctx, id)
	if err != nil {
		return fmt.Errorf("performing SqlResourcesMigrateSqlContainerToManualThroughput: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SqlResourcesMigrateSqlContainerToManualThroughput: %+v", err)
	}

	return nil
}

// preparerForSqlResourcesMigrateSqlContainerToManualThroughput prepares the SqlResourcesMigrateSqlContainerToManualThroughput request.
func (c CosmosDBClient) preparerForSqlResourcesMigrateSqlContainerToManualThroughput(ctx context.Context, id ContainerId) (*http.Request, error) {
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

// senderForSqlResourcesMigrateSqlContainerToManualThroughput sends the SqlResourcesMigrateSqlContainerToManualThroughput request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForSqlResourcesMigrateSqlContainerToManualThroughput(ctx context.Context, req *http.Request) (future SqlResourcesMigrateSqlContainerToManualThroughputOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
