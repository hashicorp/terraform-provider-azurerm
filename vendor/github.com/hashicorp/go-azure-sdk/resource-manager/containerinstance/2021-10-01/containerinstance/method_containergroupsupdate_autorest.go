package containerinstance

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerGroupsUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *ContainerGroup
}

// ContainerGroupsUpdate ...
func (c ContainerInstanceClient) ContainerGroupsUpdate(ctx context.Context, id ContainerGroupId, input Resource) (result ContainerGroupsUpdateOperationResponse, err error) {
	req, err := c.preparerForContainerGroupsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForContainerGroupsUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForContainerGroupsUpdate prepares the ContainerGroupsUpdate request.
func (c ContainerInstanceClient) preparerForContainerGroupsUpdate(ctx context.Context, id ContainerGroupId, input Resource) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForContainerGroupsUpdate handles the response to the ContainerGroupsUpdate request. The method always
// closes the http.Response Body.
func (c ContainerInstanceClient) responderForContainerGroupsUpdate(resp *http.Response) (result ContainerGroupsUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
