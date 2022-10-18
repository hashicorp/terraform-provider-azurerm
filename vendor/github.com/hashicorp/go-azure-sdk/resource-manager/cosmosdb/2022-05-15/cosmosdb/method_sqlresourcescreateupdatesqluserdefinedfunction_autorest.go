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

type SqlResourcesCreateUpdateSqlUserDefinedFunctionOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SqlResourcesCreateUpdateSqlUserDefinedFunction ...
func (c CosmosDBClient) SqlResourcesCreateUpdateSqlUserDefinedFunction(ctx context.Context, id UserDefinedFunctionId, input SqlUserDefinedFunctionCreateUpdateParameters) (result SqlResourcesCreateUpdateSqlUserDefinedFunctionOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesCreateUpdateSqlUserDefinedFunction(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesCreateUpdateSqlUserDefinedFunction", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSqlResourcesCreateUpdateSqlUserDefinedFunction(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesCreateUpdateSqlUserDefinedFunction", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SqlResourcesCreateUpdateSqlUserDefinedFunctionThenPoll performs SqlResourcesCreateUpdateSqlUserDefinedFunction then polls until it's completed
func (c CosmosDBClient) SqlResourcesCreateUpdateSqlUserDefinedFunctionThenPoll(ctx context.Context, id UserDefinedFunctionId, input SqlUserDefinedFunctionCreateUpdateParameters) error {
	result, err := c.SqlResourcesCreateUpdateSqlUserDefinedFunction(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing SqlResourcesCreateUpdateSqlUserDefinedFunction: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SqlResourcesCreateUpdateSqlUserDefinedFunction: %+v", err)
	}

	return nil
}

// preparerForSqlResourcesCreateUpdateSqlUserDefinedFunction prepares the SqlResourcesCreateUpdateSqlUserDefinedFunction request.
func (c CosmosDBClient) preparerForSqlResourcesCreateUpdateSqlUserDefinedFunction(ctx context.Context, id UserDefinedFunctionId, input SqlUserDefinedFunctionCreateUpdateParameters) (*http.Request, error) {
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

// senderForSqlResourcesCreateUpdateSqlUserDefinedFunction sends the SqlResourcesCreateUpdateSqlUserDefinedFunction request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForSqlResourcesCreateUpdateSqlUserDefinedFunction(ctx context.Context, req *http.Request) (future SqlResourcesCreateUpdateSqlUserDefinedFunctionOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
