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

type GenerateVerificationCodeOperationResponse struct {
	HttpResponse *http.Response
	Model        *VerificationCodeResponse
}

type GenerateVerificationCodeOperationOptions struct {
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

func DefaultGenerateVerificationCodeOperationOptions() GenerateVerificationCodeOperationOptions {
	return GenerateVerificationCodeOperationOptions{}
}

func (o GenerateVerificationCodeOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.IfMatch != nil {
		out["If-Match"] = *o.IfMatch
	}

	return out
}

func (o GenerateVerificationCodeOperationOptions) toQueryString() map[string]interface{} {
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

// GenerateVerificationCode ...
func (c DpsCertificateClient) GenerateVerificationCode(ctx context.Context, id CertificateId, options GenerateVerificationCodeOperationOptions) (result GenerateVerificationCodeOperationResponse, err error) {
	req, err := c.preparerForGenerateVerificationCode(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dpscertificate.DpsCertificateClient", "GenerateVerificationCode", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "dpscertificate.DpsCertificateClient", "GenerateVerificationCode", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGenerateVerificationCode(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dpscertificate.DpsCertificateClient", "GenerateVerificationCode", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGenerateVerificationCode prepares the GenerateVerificationCode request.
func (c DpsCertificateClient) preparerForGenerateVerificationCode(ctx context.Context, id CertificateId, options GenerateVerificationCodeOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/generateVerificationCode", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGenerateVerificationCode handles the response to the GenerateVerificationCode request. The method always
// closes the http.Response Body.
func (c DpsCertificateClient) responderForGenerateVerificationCode(resp *http.Response) (result GenerateVerificationCodeOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
