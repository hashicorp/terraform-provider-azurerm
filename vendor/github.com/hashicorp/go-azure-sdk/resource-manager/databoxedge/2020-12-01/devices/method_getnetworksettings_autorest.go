package devices

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetNetworkSettingsOperationResponse struct {
	HttpResponse *http.Response
	Model        *NetworkSettings
}

// GetNetworkSettings ...
func (c DevicesClient) GetNetworkSettings(ctx context.Context, id DataBoxEdgeDeviceId) (result GetNetworkSettingsOperationResponse, err error) {
	req, err := c.preparerForGetNetworkSettings(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "GetNetworkSettings", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "GetNetworkSettings", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetNetworkSettings(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "GetNetworkSettings", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetNetworkSettings prepares the GetNetworkSettings request.
func (c DevicesClient) preparerForGetNetworkSettings(ctx context.Context, id DataBoxEdgeDeviceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/networkSettings/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetNetworkSettings handles the response to the GetNetworkSettings request. The method always
// closes the http.Response Body.
func (c DevicesClient) responderForGetNetworkSettings(resp *http.Response) (result GetNetworkSettingsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
