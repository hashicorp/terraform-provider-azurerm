package managedenvironments

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedCertificatesDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// ManagedCertificatesDelete ...
func (c ManagedEnvironmentsClient) ManagedCertificatesDelete(ctx context.Context, id ManagedCertificateId) (result ManagedCertificatesDeleteOperationResponse, err error) {
	req, err := c.preparerForManagedCertificatesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedCertificatesDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedCertificatesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForManagedCertificatesDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedCertificatesDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForManagedCertificatesDelete prepares the ManagedCertificatesDelete request.
func (c ManagedEnvironmentsClient) preparerForManagedCertificatesDelete(ctx context.Context, id ManagedCertificateId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForManagedCertificatesDelete handles the response to the ManagedCertificatesDelete request. The method always
// closes the http.Response Body.
func (c ManagedEnvironmentsClient) responderForManagedCertificatesDelete(resp *http.Response) (result ManagedCertificatesDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
