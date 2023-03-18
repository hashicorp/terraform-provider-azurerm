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

type ImportDataOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ImportData ...
func (c RedisClient) ImportData(ctx context.Context, id RediId, input ImportRDBParameters) (result ImportDataOperationResponse, err error) {
	req, err := c.preparerForImportData(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "ImportData", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForImportData(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "ImportData", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ImportDataThenPoll performs ImportData then polls until it's completed
func (c RedisClient) ImportDataThenPoll(ctx context.Context, id RediId, input ImportRDBParameters) error {
	result, err := c.ImportData(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ImportData: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ImportData: %+v", err)
	}

	return nil
}

// preparerForImportData prepares the ImportData request.
func (c RedisClient) preparerForImportData(ctx context.Context, id RediId, input ImportRDBParameters) (*http.Request, error) {
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

// senderForImportData sends the ImportData request. The method will close the
// http.Response Body if it receives an error.
func (c RedisClient) senderForImportData(ctx context.Context, req *http.Request) (future ImportDataOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
