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

type RestoreDefaultsOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// RestoreDefaults ...
func (c StorageTargetsClient) RestoreDefaults(ctx context.Context, id StorageTargetId) (result RestoreDefaultsOperationResponse, err error) {
	req, err := c.preparerForRestoreDefaults(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagetargets.StorageTargetsClient", "RestoreDefaults", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForRestoreDefaults(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagetargets.StorageTargetsClient", "RestoreDefaults", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// RestoreDefaultsThenPoll performs RestoreDefaults then polls until it's completed
func (c StorageTargetsClient) RestoreDefaultsThenPoll(ctx context.Context, id StorageTargetId) error {
	result, err := c.RestoreDefaults(ctx, id)
	if err != nil {
		return fmt.Errorf("performing RestoreDefaults: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after RestoreDefaults: %+v", err)
	}

	return nil
}

// preparerForRestoreDefaults prepares the RestoreDefaults request.
func (c StorageTargetsClient) preparerForRestoreDefaults(ctx context.Context, id StorageTargetId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/restoreDefaults", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForRestoreDefaults sends the RestoreDefaults request. The method will close the
// http.Response Body if it receives an error.
func (c StorageTargetsClient) senderForRestoreDefaults(ctx context.Context, req *http.Request) (future RestoreDefaultsOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
