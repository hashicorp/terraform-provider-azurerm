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

type SqlResourcesMigrateSqlContainerToAutoscaleOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SqlResourcesMigrateSqlContainerToAutoscale ...
func (c CosmosDBClient) SqlResourcesMigrateSqlContainerToAutoscale(ctx context.Context, id ContainerId) (result SqlResourcesMigrateSqlContainerToAutoscaleOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesMigrateSqlContainerToAutoscale(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesMigrateSqlContainerToAutoscale", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSqlResourcesMigrateSqlContainerToAutoscale(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesMigrateSqlContainerToAutoscale", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SqlResourcesMigrateSqlContainerToAutoscaleThenPoll performs SqlResourcesMigrateSqlContainerToAutoscale then polls until it's completed
func (c CosmosDBClient) SqlResourcesMigrateSqlContainerToAutoscaleThenPoll(ctx context.Context, id ContainerId) error {
	result, err := c.SqlResourcesMigrateSqlContainerToAutoscale(ctx, id)
	if err != nil {
		return fmt.Errorf("performing SqlResourcesMigrateSqlContainerToAutoscale: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SqlResourcesMigrateSqlContainerToAutoscale: %+v", err)
	}

	return nil
}

// preparerForSqlResourcesMigrateSqlContainerToAutoscale prepares the SqlResourcesMigrateSqlContainerToAutoscale request.
func (c CosmosDBClient) preparerForSqlResourcesMigrateSqlContainerToAutoscale(ctx context.Context, id ContainerId) (*http.Request, error) {
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

// senderForSqlResourcesMigrateSqlContainerToAutoscale sends the SqlResourcesMigrateSqlContainerToAutoscale request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForSqlResourcesMigrateSqlContainerToAutoscale(ctx context.Context, req *http.Request) (future SqlResourcesMigrateSqlContainerToAutoscaleOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
