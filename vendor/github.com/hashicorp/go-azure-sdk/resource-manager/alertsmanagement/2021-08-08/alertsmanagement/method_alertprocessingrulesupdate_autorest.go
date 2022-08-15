package alertsmanagement

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertProcessingRulesUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *AlertProcessingRule
}

// AlertProcessingRulesUpdate ...
func (c AlertsManagementClient) AlertProcessingRulesUpdate(ctx context.Context, id ActionRuleId, input PatchObject) (result AlertProcessingRulesUpdateOperationResponse, err error) {
	req, err := c.preparerForAlertProcessingRulesUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.AlertsManagementClient", "AlertProcessingRulesUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.AlertsManagementClient", "AlertProcessingRulesUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAlertProcessingRulesUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.AlertsManagementClient", "AlertProcessingRulesUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAlertProcessingRulesUpdate prepares the AlertProcessingRulesUpdate request.
func (c AlertsManagementClient) preparerForAlertProcessingRulesUpdate(ctx context.Context, id ActionRuleId, input PatchObject) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAlertProcessingRulesUpdate handles the response to the AlertProcessingRulesUpdate request. The method always
// closes the http.Response Body.
func (c AlertsManagementClient) responderForAlertProcessingRulesUpdate(resp *http.Response) (result AlertProcessingRulesUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
