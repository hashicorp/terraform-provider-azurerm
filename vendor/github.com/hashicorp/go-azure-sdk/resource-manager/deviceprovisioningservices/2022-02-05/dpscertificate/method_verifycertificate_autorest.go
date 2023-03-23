package dpscertificate

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VerifyCertificateOperationResponse struct {
	HttpResponse *http.Response
	Model        *CertificateResponse
}

type VerifyCertificateOperationOptions struct {
	CertificateCreated       *string
	CertificateHasPrivateKey *bool
	CertificateIsVerified    *bool
	CertificateLastUpdated   *string
	CertificateName          *string
	CertificateNonce         *string
	CertificatePurpose       *CertificatePurpose
	CertificateRawBytes      *string
	IfMatch                  *string
}

func DefaultVerifyCertificateOperationOptions() VerifyCertificateOperationOptions {
	return VerifyCertificateOperationOptions{}
}

func (o VerifyCertificateOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.IfMatch != nil {
		out["If-Match"] = *o.IfMatch
	}

	return out
}

func (o VerifyCertificateOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.CertificateCreated != nil {
		out["certificate.created"] = *o.CertificateCreated
	}

	if o.CertificateHasPrivateKey != nil {
		out["certificate.hasPrivateKey"] = *o.CertificateHasPrivateKey
	}

	if o.CertificateIsVerified != nil {
		out["certificate.isVerified"] = *o.CertificateIsVerified
	}

	if o.CertificateLastUpdated != nil {
		out["certificate.lastUpdated"] = *o.CertificateLastUpdated
	}

	if o.CertificateName != nil {
		out["certificate.name"] = *o.CertificateName
	}

	if o.CertificateNonce != nil {
		out["certificate.nonce"] = *o.CertificateNonce
	}

	if o.CertificatePurpose != nil {
		out["certificate.purpose"] = *o.CertificatePurpose
	}

	if o.CertificateRawBytes != nil {
		out["certificate.rawBytes"] = *o.CertificateRawBytes
	}

	return out
}

// VerifyCertificate ...
func (c DpsCertificateClient) VerifyCertificate(ctx context.Context, id CertificateId, input VerificationCodeRequest, options VerifyCertificateOperationOptions) (result VerifyCertificateOperationResponse, err error) {
	req, err := c.preparerForVerifyCertificate(ctx, id, input, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dpscertificate.DpsCertificateClient", "VerifyCertificate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "dpscertificate.DpsCertificateClient", "VerifyCertificate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForVerifyCertificate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dpscertificate.DpsCertificateClient", "VerifyCertificate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForVerifyCertificate prepares the VerifyCertificate request.
func (c DpsCertificateClient) preparerForVerifyCertificate(ctx context.Context, id CertificateId, input VerificationCodeRequest, options VerifyCertificateOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/verify", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForVerifyCertificate handles the response to the VerifyCertificate request. The method always
// closes the http.Response Body.
func (c DpsCertificateClient) responderForVerifyCertificate(resp *http.Response) (result VerifyCertificateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
