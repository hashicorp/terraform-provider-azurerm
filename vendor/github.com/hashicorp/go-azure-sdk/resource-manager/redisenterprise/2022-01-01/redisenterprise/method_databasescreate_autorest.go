package redisenterprise

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

type DatabasesCreateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DatabasesCreate ...
func (c RedisEnterpriseClient) DatabasesCreate(ctx context.Context, id DatabaseId, input Database) (result DatabasesCreateOperationResponse, err error) {
	req, err := c.preparerForDatabasesCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesCreate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDatabasesCreate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesCreate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DatabasesCreateThenPoll performs DatabasesCreate then polls until it's completed
func (c RedisEnterpriseClient) DatabasesCreateThenPoll(ctx context.Context, id DatabaseId, input Database) error {
	result, err := c.DatabasesCreate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DatabasesCreate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DatabasesCreate: %+v", err)
	}

	return nil
}

// preparerForDatabasesCreate prepares the DatabasesCreate request.
func (c RedisEnterpriseClient) preparerForDatabasesCreate(ctx context.Context, id DatabaseId, input Database) (*http.Request, error) {
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

// senderForDatabasesCreate sends the DatabasesCreate request. The method will close the
// http.Response Body if it receives an error.
func (c RedisEnterpriseClient) senderForDatabasesCreate(ctx context.Context, req *http.Request) (future DatabasesCreateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
