package workspaces

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

type ResyncKeysOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ResyncKeys ...
func (c WorkspacesClient) ResyncKeys(ctx context.Context, id WorkspaceId) (result ResyncKeysOperationResponse, err error) {
	req, err := c.preparerForResyncKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "ResyncKeys", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForResyncKeys(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "ResyncKeys", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ResyncKeysThenPoll performs ResyncKeys then polls until it's completed
func (c WorkspacesClient) ResyncKeysThenPoll(ctx context.Context, id WorkspaceId) error {
	result, err := c.ResyncKeys(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ResyncKeys: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ResyncKeys: %+v", err)
	}

	return nil
}

// preparerForResyncKeys prepares the ResyncKeys request.
func (c WorkspacesClient) preparerForResyncKeys(ctx context.Context, id WorkspaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/resyncKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForResyncKeys sends the ResyncKeys request. The method will close the
// http.Response Body if it receives an error.
func (c WorkspacesClient) senderForResyncKeys(ctx context.Context, req *http.Request) (future ResyncKeysOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
