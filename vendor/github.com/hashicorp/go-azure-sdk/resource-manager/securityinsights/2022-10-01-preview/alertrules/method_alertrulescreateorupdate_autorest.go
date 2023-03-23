package alertrules

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertRulesCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *AlertRule
}

// AlertRulesCreateOrUpdate ...
func (c AlertRulesClient) AlertRulesCreateOrUpdate(ctx context.Context, id AlertRuleId, input AlertRule) (result AlertRulesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForAlertRulesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertrules.AlertRulesClient", "AlertRulesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertrules.AlertRulesClient", "AlertRulesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAlertRulesCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertrules.AlertRulesClient", "AlertRulesCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAlertRulesCreateOrUpdate prepares the AlertRulesCreateOrUpdate request.
func (c AlertRulesClient) preparerForAlertRulesCreateOrUpdate(ctx context.Context, id AlertRuleId, input AlertRule) (*http.Request, error) {
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

// responderForAlertRulesCreateOrUpdate handles the response to the AlertRulesCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c AlertRulesClient) responderForAlertRulesCreateOrUpdate(resp *http.Response) (result AlertRulesCreateOrUpdateOperationResponse, err error) {
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
	model, err := unmarshalAlertRuleImplementation(respObj)
	if err != nil {
		return
	}
	result.Model = &model
	return
}
