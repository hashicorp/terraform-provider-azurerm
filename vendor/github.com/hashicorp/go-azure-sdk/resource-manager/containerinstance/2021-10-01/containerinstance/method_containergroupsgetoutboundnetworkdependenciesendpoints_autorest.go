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

type ContainerGroupsGetOutboundNetworkDependenciesEndpointsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]string
}

// ContainerGroupsGetOutboundNetworkDependenciesEndpoints ...
func (c ContainerInstanceClient) ContainerGroupsGetOutboundNetworkDependenciesEndpoints(ctx context.Context, id ContainerGroupId) (result ContainerGroupsGetOutboundNetworkDependenciesEndpointsOperationResponse, err error) {
	req, err := c.preparerForContainerGroupsGetOutboundNetworkDependenciesEndpoints(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsGetOutboundNetworkDependenciesEndpoints", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsGetOutboundNetworkDependenciesEndpoints", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForContainerGroupsGetOutboundNetworkDependenciesEndpoints(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsGetOutboundNetworkDependenciesEndpoints", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForContainerGroupsGetOutboundNetworkDependenciesEndpoints prepares the ContainerGroupsGetOutboundNetworkDependenciesEndpoints request.
func (c ContainerInstanceClient) preparerForContainerGroupsGetOutboundNetworkDependenciesEndpoints(ctx context.Context, id ContainerGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/outboundNetworkDependenciesEndpoints", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForContainerGroupsGetOutboundNetworkDependenciesEndpoints handles the response to the ContainerGroupsGetOutboundNetworkDependenciesEndpoints request. The method always
// closes the http.Response Body.
func (c ContainerInstanceClient) responderForContainerGroupsGetOutboundNetworkDependenciesEndpoints(resp *http.Response) (result ContainerGroupsGetOutboundNetworkDependenciesEndpointsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
