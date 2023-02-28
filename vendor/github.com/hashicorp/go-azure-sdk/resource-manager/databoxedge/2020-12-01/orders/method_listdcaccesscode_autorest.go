package orders

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListDCAccessCodeOperationResponse struct {
	HttpResponse *http.Response
	Model        *DCAccessCode
}

// ListDCAccessCode ...
func (c OrdersClient) ListDCAccessCode(ctx context.Context, id DataBoxEdgeDeviceId) (result ListDCAccessCodeOperationResponse, err error) {
	req, err := c.preparerForListDCAccessCode(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "orders.OrdersClient", "ListDCAccessCode", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "orders.OrdersClient", "ListDCAccessCode", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListDCAccessCode(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "orders.OrdersClient", "ListDCAccessCode", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListDCAccessCode prepares the ListDCAccessCode request.
func (c OrdersClient) preparerForListDCAccessCode(ctx context.Context, id DataBoxEdgeDeviceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/orders/default/listDCAccessCode", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListDCAccessCode handles the response to the ListDCAccessCode request. The method always
// closes the http.Response Body.
func (c OrdersClient) responderForListDCAccessCode(resp *http.Response) (result ListDCAccessCodeOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
