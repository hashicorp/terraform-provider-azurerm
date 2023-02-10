package certificate

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Certificate
}

type CreateOperationOptions struct {
	IfMatch     *string
	IfNoneMatch *string
}

func DefaultCreateOperationOptions() CreateOperationOptions {
	return CreateOperationOptions{}
}

func (o CreateOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.IfMatch != nil {
		out["If-Match"] = *o.IfMatch
	}

	if o.IfNoneMatch != nil {
		out["If-None-Match"] = *o.IfNoneMatch
	}

	return out
}

func (o CreateOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// Create ...
func (c CertificateClient) Create(ctx context.Context, id CertificateId, input CertificateCreateOrUpdateParameters, options CreateOperationOptions) (result CreateOperationResponse, err error) {
	req, err := c.preparerForCreate(ctx, id, input, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificate.CertificateClient", "Create", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificate.CertificateClient", "Create", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCreate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificate.CertificateClient", "Create", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCreate prepares the Create request.
func (c CertificateClient) preparerForCreate(ctx context.Context, id CertificateId, input CertificateCreateOrUpdateParameters, options CreateOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCreate handles the response to the Create request. The method always
// closes the http.Response Body.
func (c CertificateClient) responderForCreate(resp *http.Response) (result CreateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
