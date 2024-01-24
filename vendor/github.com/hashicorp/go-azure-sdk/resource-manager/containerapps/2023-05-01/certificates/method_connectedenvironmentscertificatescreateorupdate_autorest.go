package certificates

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedEnvironmentsCertificatesCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Certificate
}

// ConnectedEnvironmentsCertificatesCreateOrUpdate ...
func (c CertificatesClient) ConnectedEnvironmentsCertificatesCreateOrUpdate(ctx context.Context, id ConnectedEnvironmentCertificateId, input Certificate) (result ConnectedEnvironmentsCertificatesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForConnectedEnvironmentsCertificatesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificates.CertificatesClient", "ConnectedEnvironmentsCertificatesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificates.CertificatesClient", "ConnectedEnvironmentsCertificatesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConnectedEnvironmentsCertificatesCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificates.CertificatesClient", "ConnectedEnvironmentsCertificatesCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConnectedEnvironmentsCertificatesCreateOrUpdate prepares the ConnectedEnvironmentsCertificatesCreateOrUpdate request.
func (c CertificatesClient) preparerForConnectedEnvironmentsCertificatesCreateOrUpdate(ctx context.Context, id ConnectedEnvironmentCertificateId, input Certificate) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForConnectedEnvironmentsCertificatesCreateOrUpdate handles the response to the ConnectedEnvironmentsCertificatesCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c CertificatesClient) responderForConnectedEnvironmentsCertificatesCreateOrUpdate(resp *http.Response) (result ConnectedEnvironmentsCertificatesCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
