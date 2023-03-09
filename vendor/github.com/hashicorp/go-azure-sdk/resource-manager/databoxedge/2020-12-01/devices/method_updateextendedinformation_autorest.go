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

type UpdateExtendedInformationOperationResponse struct {
	HttpResponse *http.Response
	Model        *DataBoxEdgeDeviceExtendedInfo
}

// UpdateExtendedInformation ...
func (c DevicesClient) UpdateExtendedInformation(ctx context.Context, id DataBoxEdgeDeviceId, input DataBoxEdgeDeviceExtendedInfoPatch) (result UpdateExtendedInformationOperationResponse, err error) {
	req, err := c.preparerForUpdateExtendedInformation(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "UpdateExtendedInformation", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "UpdateExtendedInformation", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUpdateExtendedInformation(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "UpdateExtendedInformation", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUpdateExtendedInformation prepares the UpdateExtendedInformation request.
func (c DevicesClient) preparerForUpdateExtendedInformation(ctx context.Context, id DataBoxEdgeDeviceId, input DataBoxEdgeDeviceExtendedInfoPatch) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/updateExtendedInformation", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForUpdateExtendedInformation handles the response to the UpdateExtendedInformation request. The method always
// closes the http.Response Body.
func (c DevicesClient) responderForUpdateExtendedInformation(resp *http.Response) (result UpdateExtendedInformationOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
