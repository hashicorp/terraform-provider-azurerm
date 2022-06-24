package resourceguards

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetDefaultDeleteResourceGuardProxyRequestsObjectOperationResponse struct {
	HttpResponse *http.Response
	Model        *DppBaseResource
}

// GetDefaultDeleteResourceGuardProxyRequestsObject ...
func (c ResourceGuardsClient) GetDefaultDeleteResourceGuardProxyRequestsObject(ctx context.Context, id DeleteResourceGuardProxyRequestId) (result GetDefaultDeleteResourceGuardProxyRequestsObjectOperationResponse, err error) {
	req, err := c.preparerForGetDefaultDeleteResourceGuardProxyRequestsObject(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDefaultDeleteResourceGuardProxyRequestsObject", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDefaultDeleteResourceGuardProxyRequestsObject", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetDefaultDeleteResourceGuardProxyRequestsObject(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDefaultDeleteResourceGuardProxyRequestsObject", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetDefaultDeleteResourceGuardProxyRequestsObject prepares the GetDefaultDeleteResourceGuardProxyRequestsObject request.
func (c ResourceGuardsClient) preparerForGetDefaultDeleteResourceGuardProxyRequestsObject(ctx context.Context, id DeleteResourceGuardProxyRequestId) (*http.Request, error) {
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

// responderForGetDefaultDeleteResourceGuardProxyRequestsObject handles the response to the GetDefaultDeleteResourceGuardProxyRequestsObject request. The method always
// closes the http.Response Body.
func (c ResourceGuardsClient) responderForGetDefaultDeleteResourceGuardProxyRequestsObject(resp *http.Response) (result GetDefaultDeleteResourceGuardProxyRequestsObjectOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
