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

type StorageTargetInvalidateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// StorageTargetInvalidate ...
func (c StorageTargetsClient) StorageTargetInvalidate(ctx context.Context, id StorageTargetId) (result StorageTargetInvalidateOperationResponse, err error) {
	req, err := c.preparerForStorageTargetInvalidate(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagetargets.StorageTargetsClient", "StorageTargetInvalidate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForStorageTargetInvalidate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagetargets.StorageTargetsClient", "StorageTargetInvalidate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// StorageTargetInvalidateThenPoll performs StorageTargetInvalidate then polls until it's completed
func (c StorageTargetsClient) StorageTargetInvalidateThenPoll(ctx context.Context, id StorageTargetId) error {
	result, err := c.StorageTargetInvalidate(ctx, id)
	if err != nil {
		return fmt.Errorf("performing StorageTargetInvalidate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after StorageTargetInvalidate: %+v", err)
	}

	return nil
}

// preparerForStorageTargetInvalidate prepares the StorageTargetInvalidate request.
func (c StorageTargetsClient) preparerForStorageTargetInvalidate(ctx context.Context, id StorageTargetId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/invalidate", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForStorageTargetInvalidate sends the StorageTargetInvalidate request. The method will close the
// http.Response Body if it receives an error.
func (c StorageTargetsClient) senderForStorageTargetInvalidate(ctx context.Context, req *http.Request) (future StorageTargetInvalidateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
