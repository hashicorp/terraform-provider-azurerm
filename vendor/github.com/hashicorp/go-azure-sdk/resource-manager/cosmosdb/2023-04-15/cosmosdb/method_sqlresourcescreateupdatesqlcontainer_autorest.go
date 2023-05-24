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

type SqlResourcesCreateUpdateSqlContainerOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SqlResourcesCreateUpdateSqlContainer ...
func (c CosmosDBClient) SqlResourcesCreateUpdateSqlContainer(ctx context.Context, id ContainerId, input SqlContainerCreateUpdateParameters) (result SqlResourcesCreateUpdateSqlContainerOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesCreateUpdateSqlContainer(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesCreateUpdateSqlContainer", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSqlResourcesCreateUpdateSqlContainer(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesCreateUpdateSqlContainer", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SqlResourcesCreateUpdateSqlContainerThenPoll performs SqlResourcesCreateUpdateSqlContainer then polls until it's completed
func (c CosmosDBClient) SqlResourcesCreateUpdateSqlContainerThenPoll(ctx context.Context, id ContainerId, input SqlContainerCreateUpdateParameters) error {
	result, err := c.SqlResourcesCreateUpdateSqlContainer(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing SqlResourcesCreateUpdateSqlContainer: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SqlResourcesCreateUpdateSqlContainer: %+v", err)
	}

	return nil
}

// preparerForSqlResourcesCreateUpdateSqlContainer prepares the SqlResourcesCreateUpdateSqlContainer request.
func (c CosmosDBClient) preparerForSqlResourcesCreateUpdateSqlContainer(ctx context.Context, id ContainerId, input SqlContainerCreateUpdateParameters) (*http.Request, error) {
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

// senderForSqlResourcesCreateUpdateSqlContainer sends the SqlResourcesCreateUpdateSqlContainer request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForSqlResourcesCreateUpdateSqlContainer(ctx context.Context, req *http.Request) (future SqlResourcesCreateUpdateSqlContainerOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
