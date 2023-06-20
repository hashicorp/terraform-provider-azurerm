package alertrules

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertRulesDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// AlertRulesDelete ...
func (c AlertRulesClient) AlertRulesDelete(ctx context.Context, id AlertRuleId) (result AlertRulesDeleteOperationResponse, err error) {
	req, err := c.preparerForAlertRulesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertrules.AlertRulesClient", "AlertRulesDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertrules.AlertRulesClient", "AlertRulesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAlertRulesDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertrules.AlertRulesClient", "AlertRulesDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAlertRulesDelete prepares the AlertRulesDelete request.
func (c AlertRulesClient) preparerForAlertRulesDelete(ctx context.Context, id AlertRuleId) (*http.Request, error) {
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

// responderForAlertRulesDelete handles the response to the AlertRulesDelete request. The method always
// closes the http.Response Body.
func (c AlertRulesClient) responderForAlertRulesDelete(resp *http.Response) (result AlertRulesDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
