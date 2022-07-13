package containerinstance

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerGroupsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ContainerGroup
}

// ContainerGroupsGet ...
func (c ContainerInstanceClient) ContainerGroupsGet(ctx context.Context, id ContainerGroupId) (result ContainerGroupsGetOperationResponse, err error) {
	req, err := c.preparerForContainerGroupsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForContainerGroupsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForContainerGroupsGet prepares the ContainerGroupsGet request.
func (c ContainerInstanceClient) preparerForContainerGroupsGet(ctx context.Context, id ContainerGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForContainerGroupsGet handles the response to the ContainerGroupsGet request. The method always
// closes the http.Response Body.
func (c ContainerInstanceClient) responderForContainerGroupsGet(resp *http.Response) (result ContainerGroupsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
