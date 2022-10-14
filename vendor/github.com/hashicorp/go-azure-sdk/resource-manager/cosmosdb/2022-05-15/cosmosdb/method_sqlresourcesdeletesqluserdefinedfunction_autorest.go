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

type SqlResourcesDeleteSqlUserDefinedFunctionOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SqlResourcesDeleteSqlUserDefinedFunction ...
func (c CosmosDBClient) SqlResourcesDeleteSqlUserDefinedFunction(ctx context.Context, id UserDefinedFunctionId) (result SqlResourcesDeleteSqlUserDefinedFunctionOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesDeleteSqlUserDefinedFunction(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesDeleteSqlUserDefinedFunction", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSqlResourcesDeleteSqlUserDefinedFunction(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesDeleteSqlUserDefinedFunction", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SqlResourcesDeleteSqlUserDefinedFunctionThenPoll performs SqlResourcesDeleteSqlUserDefinedFunction then polls until it's completed
func (c CosmosDBClient) SqlResourcesDeleteSqlUserDefinedFunctionThenPoll(ctx context.Context, id UserDefinedFunctionId) error {
	result, err := c.SqlResourcesDeleteSqlUserDefinedFunction(ctx, id)
	if err != nil {
		return fmt.Errorf("performing SqlResourcesDeleteSqlUserDefinedFunction: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SqlResourcesDeleteSqlUserDefinedFunction: %+v", err)
	}

	return nil
}

// preparerForSqlResourcesDeleteSqlUserDefinedFunction prepares the SqlResourcesDeleteSqlUserDefinedFunction request.
func (c CosmosDBClient) preparerForSqlResourcesDeleteSqlUserDefinedFunction(ctx context.Context, id UserDefinedFunctionId) (*http.Request, error) {
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

// senderForSqlResourcesDeleteSqlUserDefinedFunction sends the SqlResourcesDeleteSqlUserDefinedFunction request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForSqlResourcesDeleteSqlUserDefinedFunction(ctx context.Context, req *http.Request) (future SqlResourcesDeleteSqlUserDefinedFunctionOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
