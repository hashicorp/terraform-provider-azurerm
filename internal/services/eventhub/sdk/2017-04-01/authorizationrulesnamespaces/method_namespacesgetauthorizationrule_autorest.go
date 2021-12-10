package authorizationrulesnamespaces

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type NamespacesGetAuthorizationRuleResponse struct {
	HttpResponse *http.Response
	Model        *AuthorizationRule
}

// NamespacesGetAuthorizationRule ...
func (c AuthorizationRulesNamespacesClient) NamespacesGetAuthorizationRule(ctx context.Context, id AuthorizationRuleId) (result NamespacesGetAuthorizationRuleResponse, err error) {
	req, err := c.preparerForNamespacesGetAuthorizationRule(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesGetAuthorizationRule", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesGetAuthorizationRule", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForNamespacesGetAuthorizationRule(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesGetAuthorizationRule", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForNamespacesGetAuthorizationRule prepares the NamespacesGetAuthorizationRule request.
func (c AuthorizationRulesNamespacesClient) preparerForNamespacesGetAuthorizationRule(ctx context.Context, id AuthorizationRuleId) (*http.Request, error) {
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

// responderForNamespacesGetAuthorizationRule handles the response to the NamespacesGetAuthorizationRule request. The method always
// closes the http.Response Body.
func (c AuthorizationRulesNamespacesClient) responderForNamespacesGetAuthorizationRule(resp *http.Response) (result NamespacesGetAuthorizationRuleResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
