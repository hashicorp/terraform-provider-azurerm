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

type TableResourcesMigrateTableToAutoscaleOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// TableResourcesMigrateTableToAutoscale ...
func (c CosmosDBClient) TableResourcesMigrateTableToAutoscale(ctx context.Context, id TableId) (result TableResourcesMigrateTableToAutoscaleOperationResponse, err error) {
	req, err := c.preparerForTableResourcesMigrateTableToAutoscale(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "TableResourcesMigrateTableToAutoscale", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForTableResourcesMigrateTableToAutoscale(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "TableResourcesMigrateTableToAutoscale", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// TableResourcesMigrateTableToAutoscaleThenPoll performs TableResourcesMigrateTableToAutoscale then polls until it's completed
func (c CosmosDBClient) TableResourcesMigrateTableToAutoscaleThenPoll(ctx context.Context, id TableId) error {
	result, err := c.TableResourcesMigrateTableToAutoscale(ctx, id)
	if err != nil {
		return fmt.Errorf("performing TableResourcesMigrateTableToAutoscale: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after TableResourcesMigrateTableToAutoscale: %+v", err)
	}

	return nil
}

// preparerForTableResourcesMigrateTableToAutoscale prepares the TableResourcesMigrateTableToAutoscale request.
func (c CosmosDBClient) preparerForTableResourcesMigrateTableToAutoscale(ctx context.Context, id TableId) (*http.Request, error) {
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

// senderForTableResourcesMigrateTableToAutoscale sends the TableResourcesMigrateTableToAutoscale request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForTableResourcesMigrateTableToAutoscale(ctx context.Context, req *http.Request) (future TableResourcesMigrateTableToAutoscaleOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
