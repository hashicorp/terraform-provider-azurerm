package authorizationrulesnamespaces

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type NamespacesCreateOrUpdateAuthorizationRuleResponse struct {
	HttpResponse *http.Response
	Model        *AuthorizationRule
}

// NamespacesCreateOrUpdateAuthorizationRule ...
func (c AuthorizationRulesNamespacesClient) NamespacesCreateOrUpdateAuthorizationRule(ctx context.Context, id AuthorizationRuleId, input AuthorizationRule) (result NamespacesCreateOrUpdateAuthorizationRuleResponse, err error) {
	req, err := c.preparerForNamespacesCreateOrUpdateAuthorizationRule(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesCreateOrUpdateAuthorizationRule", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesCreateOrUpdateAuthorizationRule", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForNamespacesCreateOrUpdateAuthorizationRule(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesCreateOrUpdateAuthorizationRule", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForNamespacesCreateOrUpdateAuthorizationRule prepares the NamespacesCreateOrUpdateAuthorizationRule request.
func (c AuthorizationRulesNamespacesClient) preparerForNamespacesCreateOrUpdateAuthorizationRule(ctx context.Context, id AuthorizationRuleId, input AuthorizationRule) (*http.Request, error) {
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

// responderForNamespacesCreateOrUpdateAuthorizationRule handles the response to the NamespacesCreateOrUpdateAuthorizationRule request. The method always
// closes the http.Response Body.
func (c AuthorizationRulesNamespacesClient) responderForNamespacesCreateOrUpdateAuthorizationRule(resp *http.Response) (result NamespacesCreateOrUpdateAuthorizationRuleResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
