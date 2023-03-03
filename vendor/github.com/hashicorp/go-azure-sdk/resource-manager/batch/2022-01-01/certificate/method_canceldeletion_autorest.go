package certificate

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CancelDeletionOperationResponse struct {
	HttpResponse *http.Response
	Model        *Certificate
}

// CancelDeletion ...
func (c CertificateClient) CancelDeletion(ctx context.Context, id CertificateId) (result CancelDeletionOperationResponse, err error) {
	req, err := c.preparerForCancelDeletion(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificate.CertificateClient", "CancelDeletion", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificate.CertificateClient", "CancelDeletion", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCancelDeletion(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificate.CertificateClient", "CancelDeletion", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCancelDeletion prepares the CancelDeletion request.
func (c CertificateClient) preparerForCancelDeletion(ctx context.Context, id CertificateId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/cancelDelete", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCancelDeletion handles the response to the CancelDeletion request. The method always
// closes the http.Response Body.
func (c CertificateClient) responderForCancelDeletion(resp *http.Response) (result CancelDeletionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
