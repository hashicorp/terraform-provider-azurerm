package caches

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

type UpgradeFirmwareOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// UpgradeFirmware ...
func (c CachesClient) UpgradeFirmware(ctx context.Context, id CacheId) (result UpgradeFirmwareOperationResponse, err error) {
	req, err := c.preparerForUpgradeFirmware(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "caches.CachesClient", "UpgradeFirmware", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForUpgradeFirmware(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "caches.CachesClient", "UpgradeFirmware", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// UpgradeFirmwareThenPoll performs UpgradeFirmware then polls until it's completed
func (c CachesClient) UpgradeFirmwareThenPoll(ctx context.Context, id CacheId) error {
	result, err := c.UpgradeFirmware(ctx, id)
	if err != nil {
		return fmt.Errorf("performing UpgradeFirmware: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after UpgradeFirmware: %+v", err)
	}

	return nil
}

// preparerForUpgradeFirmware prepares the UpgradeFirmware request.
func (c CachesClient) preparerForUpgradeFirmware(ctx context.Context, id CacheId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/upgrade", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForUpgradeFirmware sends the UpgradeFirmware request. The method will close the
// http.Response Body if it receives an error.
func (c CachesClient) senderForUpgradeFirmware(ctx context.Context, req *http.Request) (future UpgradeFirmwareOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
