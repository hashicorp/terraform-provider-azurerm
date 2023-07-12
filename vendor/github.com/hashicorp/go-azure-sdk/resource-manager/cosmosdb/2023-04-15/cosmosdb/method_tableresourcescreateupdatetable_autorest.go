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

type TableResourcesCreateUpdateTableOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// TableResourcesCreateUpdateTable ...
func (c CosmosDBClient) TableResourcesCreateUpdateTable(ctx context.Context, id TableId, input TableCreateUpdateParameters) (result TableResourcesCreateUpdateTableOperationResponse, err error) {
	req, err := c.preparerForTableResourcesCreateUpdateTable(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "TableResourcesCreateUpdateTable", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForTableResourcesCreateUpdateTable(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "TableResourcesCreateUpdateTable", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// TableResourcesCreateUpdateTableThenPoll performs TableResourcesCreateUpdateTable then polls until it's completed
func (c CosmosDBClient) TableResourcesCreateUpdateTableThenPoll(ctx context.Context, id TableId, input TableCreateUpdateParameters) error {
	result, err := c.TableResourcesCreateUpdateTable(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing TableResourcesCreateUpdateTable: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after TableResourcesCreateUpdateTable: %+v", err)
	}

	return nil
}

// preparerForTableResourcesCreateUpdateTable prepares the TableResourcesCreateUpdateTable request.
func (c CosmosDBClient) preparerForTableResourcesCreateUpdateTable(ctx context.Context, id TableId, input TableCreateUpdateParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForTableResourcesCreateUpdateTable sends the TableResourcesCreateUpdateTable request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForTableResourcesCreateUpdateTable(ctx context.Context, req *http.Request) (future TableResourcesCreateUpdateTableOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
