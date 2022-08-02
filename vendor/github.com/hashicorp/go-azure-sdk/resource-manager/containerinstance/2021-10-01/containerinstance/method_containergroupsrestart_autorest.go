package containerinstance

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

type ContainerGroupsRestartOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ContainerGroupsRestart ...
func (c ContainerInstanceClient) ContainerGroupsRestart(ctx context.Context, id ContainerGroupId) (result ContainerGroupsRestartOperationResponse, err error) {
	req, err := c.preparerForContainerGroupsRestart(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsRestart", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForContainerGroupsRestart(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsRestart", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ContainerGroupsRestartThenPoll performs ContainerGroupsRestart then polls until it's completed
func (c ContainerInstanceClient) ContainerGroupsRestartThenPoll(ctx context.Context, id ContainerGroupId) error {
	result, err := c.ContainerGroupsRestart(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ContainerGroupsRestart: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ContainerGroupsRestart: %+v", err)
	}

	return nil
}

// preparerForContainerGroupsRestart prepares the ContainerGroupsRestart request.
func (c ContainerInstanceClient) preparerForContainerGroupsRestart(ctx context.Context, id ContainerGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/restart", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForContainerGroupsRestart sends the ContainerGroupsRestart request. The method will close the
// http.Response Body if it receives an error.
func (c ContainerInstanceClient) senderForContainerGroupsRestart(ctx context.Context, req *http.Request) (future ContainerGroupsRestartOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
