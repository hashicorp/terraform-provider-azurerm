package deviceupdates

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InstancesHeadOperationResponse struct {
	HttpResponse *http.Response
}

// InstancesHead ...
func (c DeviceupdatesClient) InstancesHead(ctx context.Context, id InstanceId) (result InstancesHeadOperationResponse, err error) {
	req, err := c.preparerForInstancesHead(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "InstancesHead", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "InstancesHead", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForInstancesHead(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "InstancesHead", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForInstancesHead prepares the InstancesHead request.
func (c DeviceupdatesClient) preparerForInstancesHead(ctx context.Context, id InstanceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsHead(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForInstancesHead handles the response to the InstancesHead request. The method always
// closes the http.Response Body.
func (c DeviceupdatesClient) responderForInstancesHead(resp *http.Response) (result InstancesHeadOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
