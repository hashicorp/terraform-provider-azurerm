package wcfrelays

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type GetAuthorizationRuleResponse struct {
	HttpResponse *http.Response
	Model        *AuthorizationRule
}

// GetAuthorizationRule ...
func (c WCFRelaysClient) GetAuthorizationRule(ctx context.Context, id WcfRelayAuthorizationRuleId) (result GetAuthorizationRuleResponse, err error) {
	req, err := c.preparerForGetAuthorizationRule(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "wcfrelays.WCFRelaysClient", "GetAuthorizationRule", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "wcfrelays.WCFRelaysClient", "GetAuthorizationRule", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetAuthorizationRule(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "wcfrelays.WCFRelaysClient", "GetAuthorizationRule", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetAuthorizationRule prepares the GetAuthorizationRule request.
func (c WCFRelaysClient) preparerForGetAuthorizationRule(ctx context.Context, id WcfRelayAuthorizationRuleId) (*http.Request, error) {
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

// responderForGetAuthorizationRule handles the response to the GetAuthorizationRule request. The method always
// closes the http.Response Body.
func (c WCFRelaysClient) responderForGetAuthorizationRule(resp *http.Response) (result GetAuthorizationRuleResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
