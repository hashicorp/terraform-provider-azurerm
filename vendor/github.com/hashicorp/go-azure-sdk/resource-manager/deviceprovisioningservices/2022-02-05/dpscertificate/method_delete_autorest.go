package dpscertificate

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeleteOperationResponse struct {
	HttpResponse *http.Response
}

type DeleteOperationOptions struct {
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

func DefaultDeleteOperationOptions() DeleteOperationOptions {
	return DeleteOperationOptions{}
}

func (o DeleteOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.IfMatch != nil {
		out["If-Match"] = *o.IfMatch
	}

	return out
}

func (o DeleteOperationOptions) toQueryString() map[string]interface{} {
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

// Delete ...
func (c DpsCertificateClient) Delete(ctx context.Context, id CertificateId, options DeleteOperationOptions) (result DeleteOperationResponse, err error) {
	req, err := c.preparerForDelete(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dpscertificate.DpsCertificateClient", "Delete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "dpscertificate.DpsCertificateClient", "Delete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dpscertificate.DpsCertificateClient", "Delete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDelete prepares the Delete request.
func (c DpsCertificateClient) preparerForDelete(ctx context.Context, id CertificateId, options DeleteOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDelete handles the response to the Delete request. The method always
// closes the http.Response Body.
func (c DpsCertificateClient) responderForDelete(resp *http.Response) (result DeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
