package alertprocessingrules

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertProcessingRulesCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *AlertProcessingRule
}

// AlertProcessingRulesCreateOrUpdate ...
func (c AlertProcessingRulesClient) AlertProcessingRulesCreateOrUpdate(ctx context.Context, id ActionRuleId, input AlertProcessingRule) (result AlertProcessingRulesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForAlertProcessingRulesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertprocessingrules.AlertProcessingRulesClient", "AlertProcessingRulesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertprocessingrules.AlertProcessingRulesClient", "AlertProcessingRulesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAlertProcessingRulesCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertprocessingrules.AlertProcessingRulesClient", "AlertProcessingRulesCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAlertProcessingRulesCreateOrUpdate prepares the AlertProcessingRulesCreateOrUpdate request.
func (c AlertProcessingRulesClient) preparerForAlertProcessingRulesCreateOrUpdate(ctx context.Context, id ActionRuleId, input AlertProcessingRule) (*http.Request, error) {
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

// responderForAlertProcessingRulesCreateOrUpdate handles the response to the AlertProcessingRulesCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c AlertProcessingRulesClient) responderForAlertProcessingRulesCreateOrUpdate(resp *http.Response) (result AlertProcessingRulesCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
