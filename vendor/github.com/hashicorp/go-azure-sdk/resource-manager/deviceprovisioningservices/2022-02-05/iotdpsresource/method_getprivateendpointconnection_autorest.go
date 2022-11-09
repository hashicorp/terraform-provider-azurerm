package iotdpsresource

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetPrivateEndpointConnectionOperationResponse struct {
	HttpResponse *http.Response
	Model        *PrivateEndpointConnection
}

// GetPrivateEndpointConnection ...
func (c IotDpsResourceClient) GetPrivateEndpointConnection(ctx context.Context, id PrivateEndpointConnectionId) (result GetPrivateEndpointConnectionOperationResponse, err error) {
	req, err := c.preparerForGetPrivateEndpointConnection(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "GetPrivateEndpointConnection", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "GetPrivateEndpointConnection", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetPrivateEndpointConnection(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "GetPrivateEndpointConnection", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetPrivateEndpointConnection prepares the GetPrivateEndpointConnection request.
func (c IotDpsResourceClient) preparerForGetPrivateEndpointConnection(ctx context.Context, id PrivateEndpointConnectionId) (*http.Request, error) {
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

// responderForGetPrivateEndpointConnection handles the response to the GetPrivateEndpointConnection request. The method always
// closes the http.Response Body.
func (c IotDpsResourceClient) responderForGetPrivateEndpointConnection(resp *http.Response) (result GetPrivateEndpointConnectionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
