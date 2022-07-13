package confidentialledger

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LedgerGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ConfidentialLedger
}

// LedgerGet ...
func (c ConfidentialLedgerClient) LedgerGet(ctx context.Context, id LedgerId) (result LedgerGetOperationResponse, err error) {
	req, err := c.preparerForLedgerGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForLedgerGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForLedgerGet prepares the LedgerGet request.
func (c ConfidentialLedgerClient) preparerForLedgerGet(ctx context.Context, id LedgerId) (*http.Request, error) {
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

// responderForLedgerGet handles the response to the LedgerGet request. The method always
// closes the http.Response Body.
func (c ConfidentialLedgerClient) responderForLedgerGet(resp *http.Response) (result LedgerGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
