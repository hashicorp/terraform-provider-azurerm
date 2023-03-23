package deviceupdates

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InstancesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *Instance
}

// InstancesGet ...
func (c DeviceupdatesClient) InstancesGet(ctx context.Context, id InstanceId) (result InstancesGetOperationResponse, err error) {
	req, err := c.preparerForInstancesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "InstancesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "InstancesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForInstancesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "InstancesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForInstancesGet prepares the InstancesGet request.
func (c DeviceupdatesClient) preparerForInstancesGet(ctx context.Context, id InstanceId) (*http.Request, error) {
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

// responderForInstancesGet handles the response to the InstancesGet request. The method always
// closes the http.Response Body.
func (c DeviceupdatesClient) responderForInstancesGet(resp *http.Response) (result InstancesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
