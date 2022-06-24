package resourceguards

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetDefaultBackupSecurityPINRequestsObjectOperationResponse struct {
	HttpResponse *http.Response
	Model        *DppBaseResource
}

// GetDefaultBackupSecurityPINRequestsObject ...
func (c ResourceGuardsClient) GetDefaultBackupSecurityPINRequestsObject(ctx context.Context, id GetBackupSecurityPINRequestId) (result GetDefaultBackupSecurityPINRequestsObjectOperationResponse, err error) {
	req, err := c.preparerForGetDefaultBackupSecurityPINRequestsObject(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDefaultBackupSecurityPINRequestsObject", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDefaultBackupSecurityPINRequestsObject", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetDefaultBackupSecurityPINRequestsObject(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetDefaultBackupSecurityPINRequestsObject", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetDefaultBackupSecurityPINRequestsObject prepares the GetDefaultBackupSecurityPINRequestsObject request.
func (c ResourceGuardsClient) preparerForGetDefaultBackupSecurityPINRequestsObject(ctx context.Context, id GetBackupSecurityPINRequestId) (*http.Request, error) {
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

// responderForGetDefaultBackupSecurityPINRequestsObject handles the response to the GetDefaultBackupSecurityPINRequestsObject request. The method always
// closes the http.Response Body.
func (c ResourceGuardsClient) responderForGetDefaultBackupSecurityPINRequestsObject(resp *http.Response) (result GetDefaultBackupSecurityPINRequestsObjectOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
