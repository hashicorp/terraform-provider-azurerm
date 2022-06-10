package authorizationrulesdisasterrecoveryconfigs

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type DisasterRecoveryConfigsGetAuthorizationRuleOperationResponse struct {
	HttpResponse *http.Response
	Model        *AuthorizationRule
}

// DisasterRecoveryConfigsGetAuthorizationRule ...
func (c AuthorizationRulesDisasterRecoveryConfigsClient) DisasterRecoveryConfigsGetAuthorizationRule(ctx context.Context, id DisasterRecoveryConfigAuthorizationRuleId) (result DisasterRecoveryConfigsGetAuthorizationRuleOperationResponse, err error) {
	req, err := c.preparerForDisasterRecoveryConfigsGetAuthorizationRule(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesdisasterrecoveryconfigs.AuthorizationRulesDisasterRecoveryConfigsClient", "DisasterRecoveryConfigsGetAuthorizationRule", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesdisasterrecoveryconfigs.AuthorizationRulesDisasterRecoveryConfigsClient", "DisasterRecoveryConfigsGetAuthorizationRule", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDisasterRecoveryConfigsGetAuthorizationRule(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesdisasterrecoveryconfigs.AuthorizationRulesDisasterRecoveryConfigsClient", "DisasterRecoveryConfigsGetAuthorizationRule", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDisasterRecoveryConfigsGetAuthorizationRule prepares the DisasterRecoveryConfigsGetAuthorizationRule request.
func (c AuthorizationRulesDisasterRecoveryConfigsClient) preparerForDisasterRecoveryConfigsGetAuthorizationRule(ctx context.Context, id DisasterRecoveryConfigAuthorizationRuleId) (*http.Request, error) {
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

// responderForDisasterRecoveryConfigsGetAuthorizationRule handles the response to the DisasterRecoveryConfigsGetAuthorizationRule request. The method always
// closes the http.Response Body.
func (c AuthorizationRulesDisasterRecoveryConfigsClient) responderForDisasterRecoveryConfigsGetAuthorizationRule(resp *http.Response) (result DisasterRecoveryConfigsGetAuthorizationRuleOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
