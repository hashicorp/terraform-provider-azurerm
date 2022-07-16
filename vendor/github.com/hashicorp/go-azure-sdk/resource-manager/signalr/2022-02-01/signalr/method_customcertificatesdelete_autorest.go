package signalr

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomCertificatesDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// CustomCertificatesDelete ...
func (c SignalRClient) CustomCertificatesDelete(ctx context.Context, id CustomCertificateId) (result CustomCertificatesDeleteOperationResponse, err error) {
	req, err := c.preparerForCustomCertificatesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomCertificatesDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomCertificatesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCustomCertificatesDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomCertificatesDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCustomCertificatesDelete prepares the CustomCertificatesDelete request.
func (c SignalRClient) preparerForCustomCertificatesDelete(ctx context.Context, id CustomCertificateId) (*http.Request, error) {
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

// responderForCustomCertificatesDelete handles the response to the CustomCertificatesDelete request. The method always
// closes the http.Response Body.
func (c SignalRClient) responderForCustomCertificatesDelete(resp *http.Response) (result CustomCertificatesDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
