package serviceresource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServicesCheckStatusOperationResponse struct {
	HttpResponse *http.Response
	Model        *DataMigrationServiceStatusResponse
}

// ServicesCheckStatus ...
func (c ServiceResourceClient) ServicesCheckStatus(ctx context.Context, id ServiceId) (result ServicesCheckStatusOperationResponse, err error) {
	req, err := c.preparerForServicesCheckStatus(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "ServicesCheckStatus", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "ServicesCheckStatus", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForServicesCheckStatus(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "ServicesCheckStatus", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForServicesCheckStatus prepares the ServicesCheckStatus request.
func (c ServiceResourceClient) preparerForServicesCheckStatus(ctx context.Context, id ServiceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/checkStatus", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForServicesCheckStatus handles the response to the ServicesCheckStatus request. The method always
// closes the http.Response Body.
func (c ServiceResourceClient) responderForServicesCheckStatus(resp *http.Response) (result ServicesCheckStatusOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
