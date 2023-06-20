package sqldedicatedgateway

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ServiceResource
}

// ServiceGet ...
func (c SqlDedicatedGatewayClient) ServiceGet(ctx context.Context, id ServiceId) (result ServiceGetOperationResponse, err error) {
	req, err := c.preparerForServiceGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "sqldedicatedgateway.SqlDedicatedGatewayClient", "ServiceGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "sqldedicatedgateway.SqlDedicatedGatewayClient", "ServiceGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForServiceGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "sqldedicatedgateway.SqlDedicatedGatewayClient", "ServiceGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForServiceGet prepares the ServiceGet request.
func (c SqlDedicatedGatewayClient) preparerForServiceGet(ctx context.Context, id ServiceId) (*http.Request, error) {
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

// responderForServiceGet handles the response to the ServiceGet request. The method always
// closes the http.Response Body.
func (c SqlDedicatedGatewayClient) responderForServiceGet(resp *http.Response) (result ServiceGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
