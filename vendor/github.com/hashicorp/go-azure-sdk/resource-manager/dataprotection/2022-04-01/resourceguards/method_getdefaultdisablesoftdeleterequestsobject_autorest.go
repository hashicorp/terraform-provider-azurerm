package resourceguards

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetDefaultDisableSoftDeleteRequestsObjectOperationResponse struct {
	HttpResponse *http.Response
	Model        *DppBaseResource
}

// GetDefaultDisableSoftDeleteRequestsObject ...
func (c ResourceGuardsClient) GetDefaultDisableSoftDeleteRequestsObject(ctx context.Context, id DisableSoftDeleteRequestId) (result GetDefaultDisableSoftDeleteRequestsObjectOperationResponse, err error) {
	req, err := c.preparerForGetDefaultDisableSoftDeleteRequestsObject(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDefaultDisableSoftDeleteRequestsObject", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDefaultDisableSoftDeleteRequestsObject", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetDefaultDisableSoftDeleteRequestsObject(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDefaultDisableSoftDeleteRequestsObject", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetDefaultDisableSoftDeleteRequestsObject prepares the GetDefaultDisableSoftDeleteRequestsObject request.
func (c ResourceGuardsClient) preparerForGetDefaultDisableSoftDeleteRequestsObject(ctx context.Context, id DisableSoftDeleteRequestId) (*http.Request, error) {
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

// responderForGetDefaultDisableSoftDeleteRequestsObject handles the response to the GetDefaultDisableSoftDeleteRequestsObject request. The method always
// closes the http.Response Body.
func (c ResourceGuardsClient) responderForGetDefaultDisableSoftDeleteRequestsObject(resp *http.Response) (result GetDefaultDisableSoftDeleteRequestsObjectOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
