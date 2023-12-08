package certificates

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedEnvironmentsCertificatesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *Certificate
}

// ConnectedEnvironmentsCertificatesGet ...
func (c CertificatesClient) ConnectedEnvironmentsCertificatesGet(ctx context.Context, id ConnectedEnvironmentCertificateId) (result ConnectedEnvironmentsCertificatesGetOperationResponse, err error) {
	req, err := c.preparerForConnectedEnvironmentsCertificatesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificates.CertificatesClient", "ConnectedEnvironmentsCertificatesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificates.CertificatesClient", "ConnectedEnvironmentsCertificatesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConnectedEnvironmentsCertificatesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificates.CertificatesClient", "ConnectedEnvironmentsCertificatesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConnectedEnvironmentsCertificatesGet prepares the ConnectedEnvironmentsCertificatesGet request.
func (c CertificatesClient) preparerForConnectedEnvironmentsCertificatesGet(ctx context.Context, id ConnectedEnvironmentCertificateId) (*http.Request, error) {
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

// responderForConnectedEnvironmentsCertificatesGet handles the response to the ConnectedEnvironmentsCertificatesGet request. The method always
// closes the http.Response Body.
func (c CertificatesClient) responderForConnectedEnvironmentsCertificatesGet(resp *http.Response) (result ConnectedEnvironmentsCertificatesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
