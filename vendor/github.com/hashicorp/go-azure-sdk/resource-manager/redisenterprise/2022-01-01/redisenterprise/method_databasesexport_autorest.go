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

type DatabasesExportOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DatabasesExport ...
func (c RedisEnterpriseClient) DatabasesExport(ctx context.Context, id DatabaseId, input ExportClusterParameters) (result DatabasesExportOperationResponse, err error) {
	req, err := c.preparerForDatabasesExport(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesExport", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDatabasesExport(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesExport", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DatabasesExportThenPoll performs DatabasesExport then polls until it's completed
func (c RedisEnterpriseClient) DatabasesExportThenPoll(ctx context.Context, id DatabaseId, input ExportClusterParameters) error {
	result, err := c.DatabasesExport(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DatabasesExport: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DatabasesExport: %+v", err)
	}

	return nil
}

// preparerForDatabasesExport prepares the DatabasesExport request.
func (c RedisEnterpriseClient) preparerForDatabasesExport(ctx context.Context, id DatabaseId, input ExportClusterParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/export", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDatabasesExport sends the DatabasesExport request. The method will close the
// http.Response Body if it receives an error.
func (c RedisEnterpriseClient) senderForDatabasesExport(ctx context.Context, req *http.Request) (future DatabasesExportOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
