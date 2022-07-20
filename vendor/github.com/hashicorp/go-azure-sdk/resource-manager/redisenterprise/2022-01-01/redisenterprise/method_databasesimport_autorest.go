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

type DatabasesImportOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DatabasesImport ...
func (c RedisEnterpriseClient) DatabasesImport(ctx context.Context, id DatabaseId, input ImportClusterParameters) (result DatabasesImportOperationResponse, err error) {
	req, err := c.preparerForDatabasesImport(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesImport", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDatabasesImport(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesImport", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DatabasesImportThenPoll performs DatabasesImport then polls until it's completed
func (c RedisEnterpriseClient) DatabasesImportThenPoll(ctx context.Context, id DatabaseId, input ImportClusterParameters) error {
	result, err := c.DatabasesImport(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DatabasesImport: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DatabasesImport: %+v", err)
	}

	return nil
}

// preparerForDatabasesImport prepares the DatabasesImport request.
func (c RedisEnterpriseClient) preparerForDatabasesImport(ctx context.Context, id DatabaseId, input ImportClusterParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/import", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDatabasesImport sends the DatabasesImport request. The method will close the
// http.Response Body if it receives an error.
func (c RedisEnterpriseClient) senderForDatabasesImport(ctx context.Context, req *http.Request) (future DatabasesImportOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
