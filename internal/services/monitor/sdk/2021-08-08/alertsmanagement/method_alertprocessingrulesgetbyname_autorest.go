package alertsmanagement

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type AlertProcessingRulesGetByNameOperationResponse struct {
	HttpResponse *http.Response
	Model        *AlertProcessingRule
}

// AlertProcessingRulesGetByName ...
func (c AlertsManagementClient) AlertProcessingRulesGetByName(ctx context.Context, id ActionRuleId) (result AlertProcessingRulesGetByNameOperationResponse, err error) {
	req, err := c.preparerForAlertProcessingRulesGetByName(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.AlertsManagementClient", "AlertProcessingRulesGetByName", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.AlertsManagementClient", "AlertProcessingRulesGetByName", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAlertProcessingRulesGetByName(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.AlertsManagementClient", "AlertProcessingRulesGetByName", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAlertProcessingRulesGetByName prepares the AlertProcessingRulesGetByName request.
func (c AlertsManagementClient) preparerForAlertProcessingRulesGetByName(ctx context.Context, id ActionRuleId) (*http.Request, error) {
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
func (c AlertsManagementClient) responderForAlertProcessingRulesGetByName(resp *http.Response) (result AlertProcessingRulesGetByNameOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
