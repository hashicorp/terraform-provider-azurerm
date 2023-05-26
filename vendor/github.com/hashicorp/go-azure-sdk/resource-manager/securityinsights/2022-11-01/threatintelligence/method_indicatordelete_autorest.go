package threatintelligence

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IndicatorDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// IndicatorDelete ...
func (c ThreatIntelligenceClient) IndicatorDelete(ctx context.Context, id IndicatorId) (result IndicatorDeleteOperationResponse, err error) {
	req, err := c.preparerForIndicatorDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForIndicatorDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForIndicatorDelete prepares the IndicatorDelete request.
func (c ThreatIntelligenceClient) preparerForIndicatorDelete(ctx context.Context, id IndicatorId) (*http.Request, error) {
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

// responderForIndicatorDelete handles the response to the IndicatorDelete request. The method always
// closes the http.Response Body.
func (c ThreatIntelligenceClient) responderForIndicatorDelete(resp *http.Response) (result IndicatorDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
