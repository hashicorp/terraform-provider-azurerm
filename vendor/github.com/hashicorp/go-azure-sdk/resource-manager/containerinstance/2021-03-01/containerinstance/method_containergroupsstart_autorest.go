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

type ContainerGroupsStartOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ContainerGroupsStart ...
func (c ContainerInstanceClient) ContainerGroupsStart(ctx context.Context, id ContainerGroupId) (result ContainerGroupsStartOperationResponse, err error) {
	req, err := c.preparerForContainerGroupsStart(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsStart", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForContainerGroupsStart(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsStart", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ContainerGroupsStartThenPoll performs ContainerGroupsStart then polls until it's completed
func (c ContainerInstanceClient) ContainerGroupsStartThenPoll(ctx context.Context, id ContainerGroupId) error {
	result, err := c.ContainerGroupsStart(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ContainerGroupsStart: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ContainerGroupsStart: %+v", err)
	}

	return nil
}

// preparerForContainerGroupsStart prepares the ContainerGroupsStart request.
func (c ContainerInstanceClient) preparerForContainerGroupsStart(ctx context.Context, id ContainerGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/start", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForContainerGroupsStart sends the ContainerGroupsStart request. The method will close the
// http.Response Body if it receives an error.
func (c ContainerInstanceClient) senderForContainerGroupsStart(ctx context.Context, req *http.Request) (future ContainerGroupsStartOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
