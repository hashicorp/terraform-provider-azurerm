package managedenvironments

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedCertificatesUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *ManagedCertificate
}

// ManagedCertificatesUpdate ...
func (c ManagedEnvironmentsClient) ManagedCertificatesUpdate(ctx context.Context, id ManagedCertificateId, input ManagedCertificatePatch) (result ManagedCertificatesUpdateOperationResponse, err error) {
	req, err := c.preparerForManagedCertificatesUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedCertificatesUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedCertificatesUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForManagedCertificatesUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedCertificatesUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForManagedCertificatesUpdate prepares the ManagedCertificatesUpdate request.
func (c ManagedEnvironmentsClient) preparerForManagedCertificatesUpdate(ctx context.Context, id ManagedCertificateId, input ManagedCertificatePatch) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForManagedCertificatesUpdate handles the response to the ManagedCertificatesUpdate request. The method always
// closes the http.Response Body.
func (c ManagedEnvironmentsClient) responderForManagedCertificatesUpdate(resp *http.Response) (result ManagedCertificatesUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
