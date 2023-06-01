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

type TableResourcesDeleteTableOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// TableResourcesDeleteTable ...
func (c CosmosDBClient) TableResourcesDeleteTable(ctx context.Context, id TableId) (result TableResourcesDeleteTableOperationResponse, err error) {
	req, err := c.preparerForTableResourcesDeleteTable(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "TableResourcesDeleteTable", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForTableResourcesDeleteTable(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "TableResourcesDeleteTable", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// TableResourcesDeleteTableThenPoll performs TableResourcesDeleteTable then polls until it's completed
func (c CosmosDBClient) TableResourcesDeleteTableThenPoll(ctx context.Context, id TableId) error {
	result, err := c.TableResourcesDeleteTable(ctx, id)
	if err != nil {
		return fmt.Errorf("performing TableResourcesDeleteTable: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after TableResourcesDeleteTable: %+v", err)
	}

	return nil
}

// preparerForTableResourcesDeleteTable prepares the TableResourcesDeleteTable request.
func (c CosmosDBClient) preparerForTableResourcesDeleteTable(ctx context.Context, id TableId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForTableResourcesDeleteTable sends the TableResourcesDeleteTable request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForTableResourcesDeleteTable(ctx context.Context, req *http.Request) (future TableResourcesDeleteTableOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
