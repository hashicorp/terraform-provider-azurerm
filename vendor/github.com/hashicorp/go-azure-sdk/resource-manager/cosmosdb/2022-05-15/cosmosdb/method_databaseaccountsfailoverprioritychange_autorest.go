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

type DatabaseAccountsFailoverPriorityChangeOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DatabaseAccountsFailoverPriorityChange ...
func (c CosmosDBClient) DatabaseAccountsFailoverPriorityChange(ctx context.Context, id DatabaseAccountId, input FailoverPolicies) (result DatabaseAccountsFailoverPriorityChangeOperationResponse, err error) {
	req, err := c.preparerForDatabaseAccountsFailoverPriorityChange(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsFailoverPriorityChange", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDatabaseAccountsFailoverPriorityChange(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsFailoverPriorityChange", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DatabaseAccountsFailoverPriorityChangeThenPoll performs DatabaseAccountsFailoverPriorityChange then polls until it's completed
func (c CosmosDBClient) DatabaseAccountsFailoverPriorityChangeThenPoll(ctx context.Context, id DatabaseAccountId, input FailoverPolicies) error {
	result, err := c.DatabaseAccountsFailoverPriorityChange(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DatabaseAccountsFailoverPriorityChange: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DatabaseAccountsFailoverPriorityChange: %+v", err)
	}

	return nil
}

// preparerForDatabaseAccountsFailoverPriorityChange prepares the DatabaseAccountsFailoverPriorityChange request.
func (c CosmosDBClient) preparerForDatabaseAccountsFailoverPriorityChange(ctx context.Context, id DatabaseAccountId, input FailoverPolicies) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/failoverPriorityChange", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDatabaseAccountsFailoverPriorityChange sends the DatabaseAccountsFailoverPriorityChange request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForDatabaseAccountsFailoverPriorityChange(ctx context.Context, req *http.Request) (future DatabaseAccountsFailoverPriorityChangeOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
