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

type ContainerGroupsDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ContainerGroupsDelete ...
func (c ContainerInstanceClient) ContainerGroupsDelete(ctx context.Context, id ContainerGroupId) (result ContainerGroupsDeleteOperationResponse, err error) {
	req, err := c.preparerForContainerGroupsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForContainerGroupsDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ContainerGroupsDeleteThenPoll performs ContainerGroupsDelete then polls until it's completed
func (c ContainerInstanceClient) ContainerGroupsDeleteThenPoll(ctx context.Context, id ContainerGroupId) error {
	result, err := c.ContainerGroupsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ContainerGroupsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ContainerGroupsDelete: %+v", err)
	}

	return nil
}

// preparerForContainerGroupsDelete prepares the ContainerGroupsDelete request.
func (c ContainerInstanceClient) preparerForContainerGroupsDelete(ctx context.Context, id ContainerGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForContainerGroupsDelete sends the ContainerGroupsDelete request. The method will close the
// http.Response Body if it receives an error.
func (c ContainerInstanceClient) senderForContainerGroupsDelete(ctx context.Context, req *http.Request) (future ContainerGroupsDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
