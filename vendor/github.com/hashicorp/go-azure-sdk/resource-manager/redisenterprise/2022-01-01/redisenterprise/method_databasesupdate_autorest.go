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

type DatabasesUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DatabasesUpdate ...
func (c RedisEnterpriseClient) DatabasesUpdate(ctx context.Context, id DatabaseId, input DatabaseUpdate) (result DatabasesUpdateOperationResponse, err error) {
	req, err := c.preparerForDatabasesUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDatabasesUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DatabasesUpdateThenPoll performs DatabasesUpdate then polls until it's completed
func (c RedisEnterpriseClient) DatabasesUpdateThenPoll(ctx context.Context, id DatabaseId, input DatabaseUpdate) error {
	result, err := c.DatabasesUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DatabasesUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DatabasesUpdate: %+v", err)
	}

	return nil
}

// preparerForDatabasesUpdate prepares the DatabasesUpdate request.
func (c RedisEnterpriseClient) preparerForDatabasesUpdate(ctx context.Context, id DatabaseId, input DatabaseUpdate) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDatabasesUpdate sends the DatabasesUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c RedisEnterpriseClient) senderForDatabasesUpdate(ctx context.Context, req *http.Request) (future DatabasesUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
