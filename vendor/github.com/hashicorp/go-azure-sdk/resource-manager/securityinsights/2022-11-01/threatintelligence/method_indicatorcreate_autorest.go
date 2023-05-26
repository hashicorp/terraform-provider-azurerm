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

type IndicatorCreateOperationResponse struct {
	HttpResponse *http.Response
	Model        *ThreatIntelligenceInformation
}

// IndicatorCreate ...
func (c ThreatIntelligenceClient) IndicatorCreate(ctx context.Context, id IndicatorId, input ThreatIntelligenceIndicatorModel) (result IndicatorCreateOperationResponse, err error) {
	req, err := c.preparerForIndicatorCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorCreate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorCreate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForIndicatorCreate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorCreate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForIndicatorCreate prepares the IndicatorCreate request.
func (c ThreatIntelligenceClient) preparerForIndicatorCreate(ctx context.Context, id IndicatorId, input ThreatIntelligenceIndicatorModel) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForIndicatorCreate handles the response to the IndicatorCreate request. The method always
// closes the http.Response Body.
func (c ThreatIntelligenceClient) responderForIndicatorCreate(resp *http.Response) (result IndicatorCreateOperationResponse, err error) {
	var respObj json.RawMessage
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
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
