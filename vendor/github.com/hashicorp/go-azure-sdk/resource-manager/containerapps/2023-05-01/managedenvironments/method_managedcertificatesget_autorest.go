package managedenvironments

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedCertificatesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ManagedCertificate
}

// ManagedCertificatesGet ...
func (c ManagedEnvironmentsClient) ManagedCertificatesGet(ctx context.Context, id ManagedCertificateId) (result ManagedCertificatesGetOperationResponse, err error) {
	req, err := c.preparerForManagedCertificatesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedCertificatesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedCertificatesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForManagedCertificatesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedCertificatesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForManagedCertificatesGet prepares the ManagedCertificatesGet request.
func (c ManagedEnvironmentsClient) preparerForManagedCertificatesGet(ctx context.Context, id ManagedCertificateId) (*http.Request, error) {
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

// responderForManagedCertificatesGet handles the response to the ManagedCertificatesGet request. The method always
// closes the http.Response Body.
func (c ManagedEnvironmentsClient) responderForManagedCertificatesGet(resp *http.Response) (result ManagedCertificatesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
