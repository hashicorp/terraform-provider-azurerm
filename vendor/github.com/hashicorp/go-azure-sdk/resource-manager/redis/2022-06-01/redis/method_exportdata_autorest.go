package redis

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

type ExportDataOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ExportData ...
func (c RedisClient) ExportData(ctx context.Context, id RediId, input ExportRDBParameters) (result ExportDataOperationResponse, err error) {
	req, err := c.preparerForExportData(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "ExportData", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForExportData(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "ExportData", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ExportDataThenPoll performs ExportData then polls until it's completed
func (c RedisClient) ExportDataThenPoll(ctx context.Context, id RediId, input ExportRDBParameters) error {
	result, err := c.ExportData(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ExportData: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ExportData: %+v", err)
	}

	return nil
}

// preparerForExportData prepares the ExportData request.
func (c RedisClient) preparerForExportData(ctx context.Context, id RediId, input ExportRDBParameters) (*http.Request, error) {
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

// senderForExportData sends the ExportData request. The method will close the
// http.Response Body if it receives an error.
func (c RedisClient) senderForExportData(ctx context.Context, req *http.Request) (future ExportDataOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
