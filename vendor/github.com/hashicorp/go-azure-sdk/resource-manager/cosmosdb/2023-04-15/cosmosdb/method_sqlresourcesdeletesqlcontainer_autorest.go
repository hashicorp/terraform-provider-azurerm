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

type SqlResourcesDeleteSqlContainerOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SqlResourcesDeleteSqlContainer ...
func (c CosmosDBClient) SqlResourcesDeleteSqlContainer(ctx context.Context, id ContainerId) (result SqlResourcesDeleteSqlContainerOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesDeleteSqlContainer(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesDeleteSqlContainer", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSqlResourcesDeleteSqlContainer(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesDeleteSqlContainer", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SqlResourcesDeleteSqlContainerThenPoll performs SqlResourcesDeleteSqlContainer then polls until it's completed
func (c CosmosDBClient) SqlResourcesDeleteSqlContainerThenPoll(ctx context.Context, id ContainerId) error {
	result, err := c.SqlResourcesDeleteSqlContainer(ctx, id)
	if err != nil {
		return fmt.Errorf("performing SqlResourcesDeleteSqlContainer: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SqlResourcesDeleteSqlContainer: %+v", err)
	}

	return nil
}

// preparerForSqlResourcesDeleteSqlContainer prepares the SqlResourcesDeleteSqlContainer request.
func (c CosmosDBClient) preparerForSqlResourcesDeleteSqlContainer(ctx context.Context, id ContainerId) (*http.Request, error) {
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

// senderForSqlResourcesDeleteSqlContainer sends the SqlResourcesDeleteSqlContainer request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForSqlResourcesDeleteSqlContainer(ctx context.Context, req *http.Request) (future SqlResourcesDeleteSqlContainerOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
