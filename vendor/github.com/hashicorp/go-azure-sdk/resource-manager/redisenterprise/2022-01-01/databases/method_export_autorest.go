package databases

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

type ExportOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Export ...
func (c DatabasesClient) Export(ctx context.Context, id DatabaseId, input ExportClusterParameters) (result ExportOperationResponse, err error) {
	req, err := c.preparerForExport(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "Export", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForExport(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "Export", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ExportThenPoll performs Export then polls until it's completed
func (c DatabasesClient) ExportThenPoll(ctx context.Context, id DatabaseId, input ExportClusterParameters) error {
	result, err := c.Export(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Export: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Export: %+v", err)
	}

	return nil
}

// preparerForExport prepares the Export request.
func (c DatabasesClient) preparerForExport(ctx context.Context, id DatabaseId, input ExportClusterParameters) (*http.Request, error) {
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

// senderForExport sends the Export request. The method will close the
// http.Response Body if it receives an error.
func (c DatabasesClient) senderForExport(ctx context.Context, req *http.Request) (future ExportOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
