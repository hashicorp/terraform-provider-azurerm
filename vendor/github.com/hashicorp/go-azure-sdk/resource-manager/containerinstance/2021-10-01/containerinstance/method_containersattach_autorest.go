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

type ContainersAttachOperationResponse struct {
	HttpResponse *http.Response
	Model        *ContainerAttachResponse
}

// ContainersAttach ...
func (c ContainerInstanceClient) ContainersAttach(ctx context.Context, id ContainerId) (result ContainersAttachOperationResponse, err error) {
	req, err := c.preparerForContainersAttach(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainersAttach", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainersAttach", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForContainersAttach(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainersAttach", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForContainersAttach prepares the ContainersAttach request.
func (c ContainerInstanceClient) preparerForContainersAttach(ctx context.Context, id ContainerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/attach", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForContainersAttach handles the response to the ContainersAttach request. The method always
// closes the http.Response Body.
func (c ContainerInstanceClient) responderForContainersAttach(resp *http.Response) (result ContainersAttachOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
