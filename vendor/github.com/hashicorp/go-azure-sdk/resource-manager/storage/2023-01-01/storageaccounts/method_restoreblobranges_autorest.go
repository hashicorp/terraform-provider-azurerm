package storageaccounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestoreBlobRangesOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// RestoreBlobRanges ...
func (c StorageAccountsClient) RestoreBlobRanges(ctx context.Context, id commonids.StorageAccountId, input BlobRestoreParameters) (result RestoreBlobRangesOperationResponse, err error) {
	req, err := c.preparerForRestoreBlobRanges(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "RestoreBlobRanges", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForRestoreBlobRanges(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "RestoreBlobRanges", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// RestoreBlobRangesThenPoll performs RestoreBlobRanges then polls until it's completed
func (c StorageAccountsClient) RestoreBlobRangesThenPoll(ctx context.Context, id commonids.StorageAccountId, input BlobRestoreParameters) error {
	result, err := c.RestoreBlobRanges(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing RestoreBlobRanges: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after RestoreBlobRanges: %+v", err)
	}

	return nil
}

// preparerForRestoreBlobRanges prepares the RestoreBlobRanges request.
func (c StorageAccountsClient) preparerForRestoreBlobRanges(ctx context.Context, id commonids.StorageAccountId, input BlobRestoreParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/restoreBlobRanges", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForRestoreBlobRanges sends the RestoreBlobRanges request. The method will close the
// http.Response Body if it receives an error.
func (c StorageAccountsClient) senderForRestoreBlobRanges(ctx context.Context, req *http.Request) (future RestoreBlobRangesOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
