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

type ContainersExecuteCommandOperationResponse struct {
	HttpResponse *http.Response
	Model        *ContainerExecResponse
}

// ContainersExecuteCommand ...
func (c ContainerInstanceClient) ContainersExecuteCommand(ctx context.Context, id ContainerId, input ContainerExecRequest) (result ContainersExecuteCommandOperationResponse, err error) {
	req, err := c.preparerForContainersExecuteCommand(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainersExecuteCommand", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainersExecuteCommand", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForContainersExecuteCommand(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainersExecuteCommand", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForContainersExecuteCommand prepares the ContainersExecuteCommand request.
func (c ContainerInstanceClient) preparerForContainersExecuteCommand(ctx context.Context, id ContainerId, input ContainerExecRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/exec", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForContainersExecuteCommand handles the response to the ContainersExecuteCommand request. The method always
// closes the http.Response Body.
func (c ContainerInstanceClient) responderForContainersExecuteCommand(resp *http.Response) (result ContainersExecuteCommandOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
