package resourceguards

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetDefaultDeleteProtectedItemRequestsObjectOperationResponse struct {
	HttpResponse *http.Response
	Model        *DppBaseResource
}

// GetDefaultDeleteProtectedItemRequestsObject ...
func (c ResourceGuardsClient) GetDefaultDeleteProtectedItemRequestsObject(ctx context.Context, id DeleteProtectedItemRequestId) (result GetDefaultDeleteProtectedItemRequestsObjectOperationResponse, err error) {
	req, err := c.preparerForGetDefaultDeleteProtectedItemRequestsObject(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDefaultDeleteProtectedItemRequestsObject", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDefaultDeleteProtectedItemRequestsObject", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetDefaultDeleteProtectedItemRequestsObject(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDefaultDeleteProtectedItemRequestsObject", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetDefaultDeleteProtectedItemRequestsObject prepares the GetDefaultDeleteProtectedItemRequestsObject request.
func (c ResourceGuardsClient) preparerForGetDefaultDeleteProtectedItemRequestsObject(ctx context.Context, id DeleteProtectedItemRequestId) (*http.Request, error) {
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

// responderForGetDefaultDeleteProtectedItemRequestsObject handles the response to the GetDefaultDeleteProtectedItemRequestsObject request. The method always
// closes the http.Response Body.
func (c ResourceGuardsClient) responderForGetDefaultDeleteProtectedItemRequestsObject(resp *http.Response) (result GetDefaultDeleteProtectedItemRequestsObjectOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
