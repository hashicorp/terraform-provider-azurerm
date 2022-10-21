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

type ContainerGroupsCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ContainerGroupsCreateOrUpdate ...
func (c ContainerInstanceClient) ContainerGroupsCreateOrUpdate(ctx context.Context, id ContainerGroupId, input ContainerGroup) (result ContainerGroupsCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForContainerGroupsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForContainerGroupsCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ContainerGroupsCreateOrUpdateThenPoll performs ContainerGroupsCreateOrUpdate then polls until it's completed
func (c ContainerInstanceClient) ContainerGroupsCreateOrUpdateThenPoll(ctx context.Context, id ContainerGroupId, input ContainerGroup) error {
	result, err := c.ContainerGroupsCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ContainerGroupsCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ContainerGroupsCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForContainerGroupsCreateOrUpdate prepares the ContainerGroupsCreateOrUpdate request.
func (c ContainerInstanceClient) preparerForContainerGroupsCreateOrUpdate(ctx context.Context, id ContainerGroupId, input ContainerGroup) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForContainerGroupsCreateOrUpdate sends the ContainerGroupsCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c ContainerInstanceClient) senderForContainerGroupsCreateOrUpdate(ctx context.Context, req *http.Request) (future ContainerGroupsCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
