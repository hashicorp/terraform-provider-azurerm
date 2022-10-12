package nginxcertificate

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificatesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *NginxCertificate
}

// CertificatesGet ...
func (c NginxCertificateClient) CertificatesGet(ctx context.Context, id CertificateId) (result CertificatesGetOperationResponse, err error) {
	req, err := c.preparerForCertificatesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxcertificate.NginxCertificateClient", "CertificatesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxcertificate.NginxCertificateClient", "CertificatesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCertificatesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxcertificate.NginxCertificateClient", "CertificatesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCertificatesGet prepares the CertificatesGet request.
func (c NginxCertificateClient) preparerForCertificatesGet(ctx context.Context, id CertificateId) (*http.Request, error) {
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

// responderForCertificatesGet handles the response to the CertificatesGet request. The method always
// closes the http.Response Body.
func (c NginxCertificateClient) responderForCertificatesGet(resp *http.Response) (result CertificatesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
