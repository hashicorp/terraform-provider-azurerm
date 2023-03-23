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

type DatabasesDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DatabasesDelete ...
func (c RedisEnterpriseClient) DatabasesDelete(ctx context.Context, id DatabaseId) (result DatabasesDeleteOperationResponse, err error) {
	req, err := c.preparerForDatabasesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDatabasesDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DatabasesDeleteThenPoll performs DatabasesDelete then polls until it's completed
func (c RedisEnterpriseClient) DatabasesDeleteThenPoll(ctx context.Context, id DatabaseId) error {
	result, err := c.DatabasesDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing DatabasesDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DatabasesDelete: %+v", err)
	}

	return nil
}

// preparerForDatabasesDelete prepares the DatabasesDelete request.
func (c RedisEnterpriseClient) preparerForDatabasesDelete(ctx context.Context, id DatabaseId) (*http.Request, error) {
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

// senderForDatabasesDelete sends the DatabasesDelete request. The method will close the
// http.Response Body if it receives an error.
func (c RedisEnterpriseClient) senderForDatabasesDelete(ctx context.Context, req *http.Request) (future DatabasesDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
