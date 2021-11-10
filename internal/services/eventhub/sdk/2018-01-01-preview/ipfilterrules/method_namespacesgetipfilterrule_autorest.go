package ipfilterrules

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type NamespacesGetIpFilterRuleResponse struct {
	HttpResponse *http.Response
	Model        *IpFilterRule
}

// NamespacesGetIpFilterRule ...
func (c IpFilterRulesClient) NamespacesGetIpFilterRule(ctx context.Context, id IpfilterruleId) (result NamespacesGetIpFilterRuleResponse, err error) {
	req, err := c.preparerForNamespacesGetIpFilterRule(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "ipfilterrules.IpFilterRulesClient", "NamespacesGetIpFilterRule", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "ipfilterrules.IpFilterRulesClient", "NamespacesGetIpFilterRule", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForNamespacesGetIpFilterRule(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "ipfilterrules.IpFilterRulesClient", "NamespacesGetIpFilterRule", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForNamespacesGetIpFilterRule prepares the NamespacesGetIpFilterRule request.
func (c IpFilterRulesClient) preparerForNamespacesGetIpFilterRule(ctx context.Context, id IpfilterruleId) (*http.Request, error) {
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

// responderForNamespacesGetIpFilterRule handles the response to the NamespacesGetIpFilterRule request. The method always
// closes the http.Response Body.
func (c IpFilterRulesClient) responderForNamespacesGetIpFilterRule(resp *http.Response) (result NamespacesGetIpFilterRuleResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
