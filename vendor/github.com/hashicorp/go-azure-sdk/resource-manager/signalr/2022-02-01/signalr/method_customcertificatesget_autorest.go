package signalr

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomCertificatesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *CustomCertificate
}

// CustomCertificatesGet ...
func (c SignalRClient) CustomCertificatesGet(ctx context.Context, id CustomCertificateId) (result CustomCertificatesGetOperationResponse, err error) {
	req, err := c.preparerForCustomCertificatesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomCertificatesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomCertificatesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCustomCertificatesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomCertificatesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCustomCertificatesGet prepares the CustomCertificatesGet request.
func (c SignalRClient) preparerForCustomCertificatesGet(ctx context.Context, id CustomCertificateId) (*http.Request, error) {
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

// responderForCustomCertificatesGet handles the response to the CustomCertificatesGet request. The method always
// closes the http.Response Body.
func (c SignalRClient) responderForCustomCertificatesGet(resp *http.Response) (result CustomCertificatesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
