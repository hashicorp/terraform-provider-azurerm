package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceListOperationResponse struct {
	HttpResponse *http.Response
	Model        *ServiceResourceListResult
}

// ServiceList ...
func (c ServicesClient) ServiceList(ctx context.Context, id DatabaseAccountId) (result ServiceListOperationResponse, err error) {
	req, err := c.preparerForServiceList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "services.ServicesClient", "ServiceList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "services.ServicesClient", "ServiceList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForServiceList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "services.ServicesClient", "ServiceList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForServiceList prepares the ServiceList request.
func (c ServicesClient) preparerForServiceList(ctx context.Context, id DatabaseAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/services", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForServiceList handles the response to the ServiceList request. The method always
// closes the http.Response Body.
func (c ServicesClient) responderForServiceList(resp *http.Response) (result ServiceListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
