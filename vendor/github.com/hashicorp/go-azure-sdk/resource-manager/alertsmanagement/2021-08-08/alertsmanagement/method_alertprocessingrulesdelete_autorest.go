package alertsmanagement

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertProcessingRulesDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// AlertProcessingRulesDelete ...
func (c AlertsManagementClient) AlertProcessingRulesDelete(ctx context.Context, id ActionRuleId) (result AlertProcessingRulesDeleteOperationResponse, err error) {
	req, err := c.preparerForAlertProcessingRulesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.AlertsManagementClient", "AlertProcessingRulesDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.AlertsManagementClient", "AlertProcessingRulesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAlertProcessingRulesDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.AlertsManagementClient", "AlertProcessingRulesDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAlertProcessingRulesDelete prepares the AlertProcessingRulesDelete request.
func (c AlertsManagementClient) preparerForAlertProcessingRulesDelete(ctx context.Context, id ActionRuleId) (*http.Request, error) {
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

// responderForAlertProcessingRulesDelete handles the response to the AlertProcessingRulesDelete request. The method always
// closes the http.Response Body.
func (c AlertsManagementClient) responderForAlertProcessingRulesDelete(resp *http.Response) (result AlertProcessingRulesDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
