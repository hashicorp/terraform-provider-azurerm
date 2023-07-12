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

type SqlResourcesMigrateSqlDatabaseToManualThroughputOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SqlResourcesMigrateSqlDatabaseToManualThroughput ...
func (c CosmosDBClient) SqlResourcesMigrateSqlDatabaseToManualThroughput(ctx context.Context, id SqlDatabaseId) (result SqlResourcesMigrateSqlDatabaseToManualThroughputOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesMigrateSqlDatabaseToManualThroughput(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesMigrateSqlDatabaseToManualThroughput", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSqlResourcesMigrateSqlDatabaseToManualThroughput(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesMigrateSqlDatabaseToManualThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SqlResourcesMigrateSqlDatabaseToManualThroughputThenPoll performs SqlResourcesMigrateSqlDatabaseToManualThroughput then polls until it's completed
func (c CosmosDBClient) SqlResourcesMigrateSqlDatabaseToManualThroughputThenPoll(ctx context.Context, id SqlDatabaseId) error {
	result, err := c.SqlResourcesMigrateSqlDatabaseToManualThroughput(ctx, id)
	if err != nil {
		return fmt.Errorf("performing SqlResourcesMigrateSqlDatabaseToManualThroughput: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SqlResourcesMigrateSqlDatabaseToManualThroughput: %+v", err)
	}

	return nil
}

// preparerForSqlResourcesMigrateSqlDatabaseToManualThroughput prepares the SqlResourcesMigrateSqlDatabaseToManualThroughput request.
func (c CosmosDBClient) preparerForSqlResourcesMigrateSqlDatabaseToManualThroughput(ctx context.Context, id SqlDatabaseId) (*http.Request, error) {
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

// senderForSqlResourcesMigrateSqlDatabaseToManualThroughput sends the SqlResourcesMigrateSqlDatabaseToManualThroughput request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForSqlResourcesMigrateSqlDatabaseToManualThroughput(ctx context.Context, req *http.Request) (future SqlResourcesMigrateSqlDatabaseToManualThroughputOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
