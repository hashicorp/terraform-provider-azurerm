package certificates

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedEnvironmentsCertificatesDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// ConnectedEnvironmentsCertificatesDelete ...
func (c CertificatesClient) ConnectedEnvironmentsCertificatesDelete(ctx context.Context, id ConnectedEnvironmentCertificateId) (result ConnectedEnvironmentsCertificatesDeleteOperationResponse, err error) {
	req, err := c.preparerForConnectedEnvironmentsCertificatesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificates.CertificatesClient", "ConnectedEnvironmentsCertificatesDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificates.CertificatesClient", "ConnectedEnvironmentsCertificatesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConnectedEnvironmentsCertificatesDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificates.CertificatesClient", "ConnectedEnvironmentsCertificatesDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConnectedEnvironmentsCertificatesDelete prepares the ConnectedEnvironmentsCertificatesDelete request.
func (c CertificatesClient) preparerForConnectedEnvironmentsCertificatesDelete(ctx context.Context, id ConnectedEnvironmentCertificateId) (*http.Request, error) {
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

// responderForConnectedEnvironmentsCertificatesDelete handles the response to the ConnectedEnvironmentsCertificatesDelete request. The method always
// closes the http.Response Body.
func (c CertificatesClient) responderForConnectedEnvironmentsCertificatesDelete(resp *http.Response) (result ConnectedEnvironmentsCertificatesDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
