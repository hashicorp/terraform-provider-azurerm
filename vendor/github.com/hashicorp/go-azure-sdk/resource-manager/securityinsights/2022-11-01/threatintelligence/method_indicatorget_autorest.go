package threatintelligence

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IndicatorGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ThreatIntelligenceInformation
}

// IndicatorGet ...
func (c ThreatIntelligenceClient) IndicatorGet(ctx context.Context, id IndicatorId) (result IndicatorGetOperationResponse, err error) {
	req, err := c.preparerForIndicatorGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForIndicatorGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForIndicatorGet prepares the IndicatorGet request.
func (c ThreatIntelligenceClient) preparerForIndicatorGet(ctx context.Context, id IndicatorId) (*http.Request, error) {
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

// responderForIndicatorGet handles the response to the IndicatorGet request. The method always
// closes the http.Response Body.
func (c ThreatIntelligenceClient) responderForIndicatorGet(resp *http.Response) (result IndicatorGetOperationResponse, err error) {
	var respObj json.RawMessage
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&respObj),
		autorest.ByClosing())
	result.HttpResponse = resp
	if err != nil {
		return
	}
	model, err := unmarshalThreatIntelligenceInformationImplementation(respObj)
	if err != nil {
		return
	}
	result.Model = &model
	return
}
