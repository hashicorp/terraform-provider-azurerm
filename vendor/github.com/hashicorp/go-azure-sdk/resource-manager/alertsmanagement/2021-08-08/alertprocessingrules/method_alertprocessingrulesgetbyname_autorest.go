package alertprocessingrules

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertProcessingRulesGetByNameOperationResponse struct {
	HttpResponse *http.Response
	Model        *AlertProcessingRule
}

// AlertProcessingRulesGetByName ...
func (c AlertProcessingRulesClient) AlertProcessingRulesGetByName(ctx context.Context, id ActionRuleId) (result AlertProcessingRulesGetByNameOperationResponse, err error) {
	req, err := c.preparerForAlertProcessingRulesGetByName(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertprocessingrules.AlertProcessingRulesClient", "AlertProcessingRulesGetByName", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertprocessingrules.AlertProcessingRulesClient", "AlertProcessingRulesGetByName", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAlertProcessingRulesGetByName(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertprocessingrules.AlertProcessingRulesClient", "AlertProcessingRulesGetByName", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAlertProcessingRulesGetByName prepares the AlertProcessingRulesGetByName request.
func (c AlertProcessingRulesClient) preparerForAlertProcessingRulesGetByName(ctx context.Context, id ActionRuleId) (*http.Request, error) {
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

// responderForAlertProcessingRulesGetByName handles the response to the AlertProcessingRulesGetByName request. The method always
// closes the http.Response Body.
func (c AlertProcessingRulesClient) responderForAlertProcessingRulesGetByName(resp *http.Response) (result AlertProcessingRulesGetByNameOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
