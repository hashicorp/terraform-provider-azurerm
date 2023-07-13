package extensions

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

type ExtensionsUpgradeOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ExtensionsUpgrade ...
func (c ExtensionsClient) ExtensionsUpgrade(ctx context.Context, id ExtensionId, input ExtensionUpgradeParameters) (result ExtensionsUpgradeOperationResponse, err error) {
	req, err := c.preparerForExtensionsUpgrade(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "extensions.ExtensionsClient", "ExtensionsUpgrade", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForExtensionsUpgrade(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "extensions.ExtensionsClient", "ExtensionsUpgrade", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ExtensionsUpgradeThenPoll performs ExtensionsUpgrade then polls until it's completed
func (c ExtensionsClient) ExtensionsUpgradeThenPoll(ctx context.Context, id ExtensionId, input ExtensionUpgradeParameters) error {
	result, err := c.ExtensionsUpgrade(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ExtensionsUpgrade: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ExtensionsUpgrade: %+v", err)
	}

	return nil
}

// preparerForExtensionsUpgrade prepares the ExtensionsUpgrade request.
func (c ExtensionsClient) preparerForExtensionsUpgrade(ctx context.Context, id ExtensionId, input ExtensionUpgradeParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/upgrade", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForExtensionsUpgrade sends the ExtensionsUpgrade request. The method will close the
// http.Response Body if it receives an error.
func (c ExtensionsClient) senderForExtensionsUpgrade(ctx context.Context, req *http.Request) (future ExtensionsUpgradeOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
