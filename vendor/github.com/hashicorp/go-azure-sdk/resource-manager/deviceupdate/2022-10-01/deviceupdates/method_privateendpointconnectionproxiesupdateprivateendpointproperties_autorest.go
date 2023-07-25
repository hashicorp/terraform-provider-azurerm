package deviceupdates

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointConnectionProxiesUpdatePrivateEndpointPropertiesOperationResponse struct {
	HttpResponse *http.Response
}

// PrivateEndpointConnectionProxiesUpdatePrivateEndpointProperties ...
func (c DeviceupdatesClient) PrivateEndpointConnectionProxiesUpdatePrivateEndpointProperties(ctx context.Context, id PrivateEndpointConnectionProxyId, input PrivateEndpointUpdate) (result PrivateEndpointConnectionProxiesUpdatePrivateEndpointPropertiesOperationResponse, err error) {
	req, err := c.preparerForPrivateEndpointConnectionProxiesUpdatePrivateEndpointProperties(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "PrivateEndpointConnectionProxiesUpdatePrivateEndpointProperties", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "PrivateEndpointConnectionProxiesUpdatePrivateEndpointProperties", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPrivateEndpointConnectionProxiesUpdatePrivateEndpointProperties(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "PrivateEndpointConnectionProxiesUpdatePrivateEndpointProperties", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPrivateEndpointConnectionProxiesUpdatePrivateEndpointProperties prepares the PrivateEndpointConnectionProxiesUpdatePrivateEndpointProperties request.
func (c DeviceupdatesClient) preparerForPrivateEndpointConnectionProxiesUpdatePrivateEndpointProperties(ctx context.Context, id PrivateEndpointConnectionProxyId, input PrivateEndpointUpdate) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/updatePrivateEndpointProperties", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForPrivateEndpointConnectionProxiesUpdatePrivateEndpointProperties handles the response to the PrivateEndpointConnectionProxiesUpdatePrivateEndpointProperties request. The method always
// closes the http.Response Body.
func (c DeviceupdatesClient) responderForPrivateEndpointConnectionProxiesUpdatePrivateEndpointProperties(resp *http.Response) (result PrivateEndpointConnectionProxiesUpdatePrivateEndpointPropertiesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
