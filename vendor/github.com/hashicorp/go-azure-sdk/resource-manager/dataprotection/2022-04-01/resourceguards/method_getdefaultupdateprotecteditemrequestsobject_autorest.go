package resourceguards

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetDefaultUpdateProtectedItemRequestsObjectOperationResponse struct {
	HttpResponse *http.Response
	Model        *DppBaseResource
}

// GetDefaultUpdateProtectedItemRequestsObject ...
func (c ResourceGuardsClient) GetDefaultUpdateProtectedItemRequestsObject(ctx context.Context, id UpdateProtectedItemRequestId) (result GetDefaultUpdateProtectedItemRequestsObjectOperationResponse, err error) {
	req, err := c.preparerForGetDefaultUpdateProtectedItemRequestsObject(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDefaultUpdateProtectedItemRequestsObject", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDefaultUpdateProtectedItemRequestsObject", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetDefaultUpdateProtectedItemRequestsObject(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDefaultUpdateProtectedItemRequestsObject", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetDefaultUpdateProtectedItemRequestsObject prepares the GetDefaultUpdateProtectedItemRequestsObject request.
func (c ResourceGuardsClient) preparerForGetDefaultUpdateProtectedItemRequestsObject(ctx context.Context, id UpdateProtectedItemRequestId) (*http.Request, error) {
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

// responderForGetDefaultUpdateProtectedItemRequestsObject handles the response to the GetDefaultUpdateProtectedItemRequestsObject request. The method always
// closes the http.Response Body.
func (c ResourceGuardsClient) responderForGetDefaultUpdateProtectedItemRequestsObject(resp *http.Response) (result GetDefaultUpdateProtectedItemRequestsObjectOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
