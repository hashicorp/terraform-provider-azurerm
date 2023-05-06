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

type AlertRulesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *AlertRule
}

// AlertRulesGet ...
func (c AlertRulesClient) AlertRulesGet(ctx context.Context, id AlertRuleId) (result AlertRulesGetOperationResponse, err error) {
	req, err := c.preparerForAlertRulesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertrules.AlertRulesClient", "AlertRulesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertrules.AlertRulesClient", "AlertRulesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAlertRulesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertrules.AlertRulesClient", "AlertRulesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAlertRulesGet prepares the AlertRulesGet request.
func (c AlertRulesClient) preparerForAlertRulesGet(ctx context.Context, id AlertRuleId) (*http.Request, error) {
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

// responderForAlertRulesGet handles the response to the AlertRulesGet request. The method always
// closes the http.Response Body.
func (c AlertRulesClient) responderForAlertRulesGet(resp *http.Response) (result AlertRulesGetOperationResponse, err error) {
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
	model, err := unmarshalAlertRuleImplementation(respObj)
	if err != nil {
		return
	}
	result.Model = &model
	return
}
