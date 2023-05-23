package storagetargets

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

type StorageTargetFlushOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// StorageTargetFlush ...
func (c StorageTargetsClient) StorageTargetFlush(ctx context.Context, id StorageTargetId) (result StorageTargetFlushOperationResponse, err error) {
	req, err := c.preparerForStorageTargetFlush(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagetargets.StorageTargetsClient", "StorageTargetFlush", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForStorageTargetFlush(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagetargets.StorageTargetsClient", "StorageTargetFlush", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// StorageTargetFlushThenPoll performs StorageTargetFlush then polls until it's completed
func (c StorageTargetsClient) StorageTargetFlushThenPoll(ctx context.Context, id StorageTargetId) error {
	result, err := c.StorageTargetFlush(ctx, id)
	if err != nil {
		return fmt.Errorf("performing StorageTargetFlush: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after StorageTargetFlush: %+v", err)
	}

	return nil
}

// preparerForStorageTargetFlush prepares the StorageTargetFlush request.
func (c StorageTargetsClient) preparerForStorageTargetFlush(ctx context.Context, id StorageTargetId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/flush", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForStorageTargetFlush sends the StorageTargetFlush request. The method will close the
// http.Response Body if it receives an error.
func (c StorageTargetsClient) senderForStorageTargetFlush(ctx context.Context, req *http.Request) (future StorageTargetFlushOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
