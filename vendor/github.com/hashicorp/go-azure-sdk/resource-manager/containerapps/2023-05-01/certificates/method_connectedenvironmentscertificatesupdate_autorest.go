package certificates

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedEnvironmentsCertificatesUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Certificate
}

// ConnectedEnvironmentsCertificatesUpdate ...
func (c CertificatesClient) ConnectedEnvironmentsCertificatesUpdate(ctx context.Context, id ConnectedEnvironmentCertificateId, input CertificatePatch) (result ConnectedEnvironmentsCertificatesUpdateOperationResponse, err error) {
	req, err := c.preparerForConnectedEnvironmentsCertificatesUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificates.CertificatesClient", "ConnectedEnvironmentsCertificatesUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificates.CertificatesClient", "ConnectedEnvironmentsCertificatesUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConnectedEnvironmentsCertificatesUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificates.CertificatesClient", "ConnectedEnvironmentsCertificatesUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConnectedEnvironmentsCertificatesUpdate prepares the ConnectedEnvironmentsCertificatesUpdate request.
func (c CertificatesClient) preparerForConnectedEnvironmentsCertificatesUpdate(ctx context.Context, id ConnectedEnvironmentCertificateId, input CertificatePatch) (*http.Request, error) {
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

// responderForConnectedEnvironmentsCertificatesUpdate handles the response to the ConnectedEnvironmentsCertificatesUpdate request. The method always
// closes the http.Response Body.
func (c CertificatesClient) responderForConnectedEnvironmentsCertificatesUpdate(resp *http.Response) (result ConnectedEnvironmentsCertificatesUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
