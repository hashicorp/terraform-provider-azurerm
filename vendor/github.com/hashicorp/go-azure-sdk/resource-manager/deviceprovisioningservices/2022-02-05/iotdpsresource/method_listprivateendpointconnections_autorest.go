package iotdpsresource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListPrivateEndpointConnectionsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]PrivateEndpointConnection
}

// ListPrivateEndpointConnections ...
func (c IotDpsResourceClient) ListPrivateEndpointConnections(ctx context.Context, id commonids.ProvisioningServiceId) (result ListPrivateEndpointConnectionsOperationResponse, err error) {
	req, err := c.preparerForListPrivateEndpointConnections(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "ListPrivateEndpointConnections", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "ListPrivateEndpointConnections", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListPrivateEndpointConnections(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "ListPrivateEndpointConnections", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListPrivateEndpointConnections prepares the ListPrivateEndpointConnections request.
func (c IotDpsResourceClient) preparerForListPrivateEndpointConnections(ctx context.Context, id commonids.ProvisioningServiceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/privateEndpointConnections", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListPrivateEndpointConnections handles the response to the ListPrivateEndpointConnections request. The method always
// closes the http.Response Body.
func (c IotDpsResourceClient) responderForListPrivateEndpointConnections(resp *http.Response) (result ListPrivateEndpointConnectionsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
