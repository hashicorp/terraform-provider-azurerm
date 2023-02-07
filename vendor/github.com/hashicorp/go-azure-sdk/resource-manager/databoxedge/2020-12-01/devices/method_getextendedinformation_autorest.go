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

type GetExtendedInformationOperationResponse struct {
	HttpResponse *http.Response
	Model        *DataBoxEdgeDeviceExtendedInfo
}

// GetExtendedInformation ...
func (c DevicesClient) GetExtendedInformation(ctx context.Context, id DataBoxEdgeDeviceId) (result GetExtendedInformationOperationResponse, err error) {
	req, err := c.preparerForGetExtendedInformation(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "GetExtendedInformation", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "GetExtendedInformation", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetExtendedInformation(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "GetExtendedInformation", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetExtendedInformation prepares the GetExtendedInformation request.
func (c DevicesClient) preparerForGetExtendedInformation(ctx context.Context, id DataBoxEdgeDeviceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/getExtendedInformation", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetExtendedInformation handles the response to the GetExtendedInformation request. The method always
// closes the http.Response Body.
func (c DevicesClient) responderForGetExtendedInformation(resp *http.Response) (result GetExtendedInformationOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
