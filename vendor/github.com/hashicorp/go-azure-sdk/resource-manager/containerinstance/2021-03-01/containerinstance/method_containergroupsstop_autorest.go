package containerinstance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerGroupsStopOperationResponse struct {
	HttpResponse *http.Response
}

// ContainerGroupsStop ...
func (c ContainerInstanceClient) ContainerGroupsStop(ctx context.Context, id ContainerGroupId) (result ContainerGroupsStopOperationResponse, err error) {
	req, err := c.preparerForContainerGroupsStop(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsStop", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsStop", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForContainerGroupsStop(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsStop", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForContainerGroupsStop prepares the ContainerGroupsStop request.
func (c ContainerInstanceClient) preparerForContainerGroupsStop(ctx context.Context, id ContainerGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/stop", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForContainerGroupsStop handles the response to the ContainerGroupsStop request. The method always
// closes the http.Response Body.
func (c ContainerInstanceClient) responderForContainerGroupsStop(resp *http.Response) (result ContainerGroupsStopOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
