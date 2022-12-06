package agentpools

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

type UpgradeNodeImageVersionOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// UpgradeNodeImageVersion ...
func (c AgentPoolsClient) UpgradeNodeImageVersion(ctx context.Context, id AgentPoolId) (result UpgradeNodeImageVersionOperationResponse, err error) {
	req, err := c.preparerForUpgradeNodeImageVersion(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "agentpools.AgentPoolsClient", "UpgradeNodeImageVersion", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForUpgradeNodeImageVersion(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "agentpools.AgentPoolsClient", "UpgradeNodeImageVersion", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// UpgradeNodeImageVersionThenPoll performs UpgradeNodeImageVersion then polls until it's completed
func (c AgentPoolsClient) UpgradeNodeImageVersionThenPoll(ctx context.Context, id AgentPoolId) error {
	result, err := c.UpgradeNodeImageVersion(ctx, id)
	if err != nil {
		return fmt.Errorf("performing UpgradeNodeImageVersion: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after UpgradeNodeImageVersion: %+v", err)
	}

	return nil
}

// preparerForUpgradeNodeImageVersion prepares the UpgradeNodeImageVersion request.
func (c AgentPoolsClient) preparerForUpgradeNodeImageVersion(ctx context.Context, id AgentPoolId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/upgradeNodeImageVersion", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForUpgradeNodeImageVersion sends the UpgradeNodeImageVersion request. The method will close the
// http.Response Body if it receives an error.
func (c AgentPoolsClient) senderForUpgradeNodeImageVersion(ctx context.Context, req *http.Request) (future UpgradeNodeImageVersionOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
