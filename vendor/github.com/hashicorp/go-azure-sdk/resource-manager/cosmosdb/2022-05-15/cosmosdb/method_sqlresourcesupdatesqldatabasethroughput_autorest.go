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

type SqlResourcesUpdateSqlDatabaseThroughputOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SqlResourcesUpdateSqlDatabaseThroughput ...
func (c CosmosDBClient) SqlResourcesUpdateSqlDatabaseThroughput(ctx context.Context, id SqlDatabaseId, input ThroughputSettingsUpdateParameters) (result SqlResourcesUpdateSqlDatabaseThroughputOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesUpdateSqlDatabaseThroughput(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesUpdateSqlDatabaseThroughput", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSqlResourcesUpdateSqlDatabaseThroughput(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesUpdateSqlDatabaseThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SqlResourcesUpdateSqlDatabaseThroughputThenPoll performs SqlResourcesUpdateSqlDatabaseThroughput then polls until it's completed
func (c CosmosDBClient) SqlResourcesUpdateSqlDatabaseThroughputThenPoll(ctx context.Context, id SqlDatabaseId, input ThroughputSettingsUpdateParameters) error {
	result, err := c.SqlResourcesUpdateSqlDatabaseThroughput(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing SqlResourcesUpdateSqlDatabaseThroughput: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SqlResourcesUpdateSqlDatabaseThroughput: %+v", err)
	}

	return nil
}

// preparerForSqlResourcesUpdateSqlDatabaseThroughput prepares the SqlResourcesUpdateSqlDatabaseThroughput request.
func (c CosmosDBClient) preparerForSqlResourcesUpdateSqlDatabaseThroughput(ctx context.Context, id SqlDatabaseId, input ThroughputSettingsUpdateParameters) (*http.Request, error) {
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

// senderForSqlResourcesUpdateSqlDatabaseThroughput sends the SqlResourcesUpdateSqlDatabaseThroughput request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForSqlResourcesUpdateSqlDatabaseThroughput(ctx context.Context, req *http.Request) (future SqlResourcesUpdateSqlDatabaseThroughputOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
