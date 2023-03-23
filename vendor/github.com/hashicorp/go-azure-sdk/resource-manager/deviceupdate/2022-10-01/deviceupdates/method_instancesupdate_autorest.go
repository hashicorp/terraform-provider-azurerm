package deviceupdates

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InstancesUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Instance
}

// InstancesUpdate ...
func (c DeviceupdatesClient) InstancesUpdate(ctx context.Context, id InstanceId, input TagUpdate) (result InstancesUpdateOperationResponse, err error) {
	req, err := c.preparerForInstancesUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "InstancesUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "InstancesUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForInstancesUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "InstancesUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForInstancesUpdate prepares the InstancesUpdate request.
func (c DeviceupdatesClient) preparerForInstancesUpdate(ctx context.Context, id InstanceId, input TagUpdate) (*http.Request, error) {
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

// responderForInstancesUpdate handles the response to the InstancesUpdate request. The method always
// closes the http.Response Body.
func (c DeviceupdatesClient) responderForInstancesUpdate(resp *http.Response) (result InstancesUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
