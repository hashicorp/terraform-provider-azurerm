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

type TableResourcesMigrateTableToManualThroughputOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// TableResourcesMigrateTableToManualThroughput ...
func (c CosmosDBClient) TableResourcesMigrateTableToManualThroughput(ctx context.Context, id TableId) (result TableResourcesMigrateTableToManualThroughputOperationResponse, err error) {
	req, err := c.preparerForTableResourcesMigrateTableToManualThroughput(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "TableResourcesMigrateTableToManualThroughput", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForTableResourcesMigrateTableToManualThroughput(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "TableResourcesMigrateTableToManualThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// TableResourcesMigrateTableToManualThroughputThenPoll performs TableResourcesMigrateTableToManualThroughput then polls until it's completed
func (c CosmosDBClient) TableResourcesMigrateTableToManualThroughputThenPoll(ctx context.Context, id TableId) error {
	result, err := c.TableResourcesMigrateTableToManualThroughput(ctx, id)
	if err != nil {
		return fmt.Errorf("performing TableResourcesMigrateTableToManualThroughput: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after TableResourcesMigrateTableToManualThroughput: %+v", err)
	}

	return nil
}

// preparerForTableResourcesMigrateTableToManualThroughput prepares the TableResourcesMigrateTableToManualThroughput request.
func (c CosmosDBClient) preparerForTableResourcesMigrateTableToManualThroughput(ctx context.Context, id TableId) (*http.Request, error) {
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

// senderForTableResourcesMigrateTableToManualThroughput sends the TableResourcesMigrateTableToManualThroughput request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForTableResourcesMigrateTableToManualThroughput(ctx context.Context, req *http.Request) (future TableResourcesMigrateTableToManualThroughputOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
