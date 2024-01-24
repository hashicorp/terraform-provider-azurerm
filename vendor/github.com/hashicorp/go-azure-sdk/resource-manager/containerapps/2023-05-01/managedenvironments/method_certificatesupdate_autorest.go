package managedenvironments

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificatesUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Certificate
}

// CertificatesUpdate ...
func (c ManagedEnvironmentsClient) CertificatesUpdate(ctx context.Context, id CertificateId, input CertificatePatch) (result CertificatesUpdateOperationResponse, err error) {
	req, err := c.preparerForCertificatesUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "CertificatesUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "CertificatesUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCertificatesUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "CertificatesUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCertificatesUpdate prepares the CertificatesUpdate request.
func (c ManagedEnvironmentsClient) preparerForCertificatesUpdate(ctx context.Context, id CertificateId, input CertificatePatch) (*http.Request, error) {
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

// responderForCertificatesUpdate handles the response to the CertificatesUpdate request. The method always
// closes the http.Response Body.
func (c ManagedEnvironmentsClient) responderForCertificatesUpdate(resp *http.Response) (result CertificatesUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
