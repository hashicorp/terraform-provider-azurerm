package ipfilterrules

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type NamespacesCreateOrUpdateIpFilterRuleResponse struct {
	HttpResponse *http.Response
	Model        *IpFilterRule
}

// NamespacesCreateOrUpdateIpFilterRule ...
func (c IpFilterRulesClient) NamespacesCreateOrUpdateIpFilterRule(ctx context.Context, id IpfilterruleId, input IpFilterRule) (result NamespacesCreateOrUpdateIpFilterRuleResponse, err error) {
	req, err := c.preparerForNamespacesCreateOrUpdateIpFilterRule(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "ipfilterrules.IpFilterRulesClient", "NamespacesCreateOrUpdateIpFilterRule", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "ipfilterrules.IpFilterRulesClient", "NamespacesCreateOrUpdateIpFilterRule", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForNamespacesCreateOrUpdateIpFilterRule(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "ipfilterrules.IpFilterRulesClient", "NamespacesCreateOrUpdateIpFilterRule", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForNamespacesCreateOrUpdateIpFilterRule prepares the NamespacesCreateOrUpdateIpFilterRule request.
func (c IpFilterRulesClient) preparerForNamespacesCreateOrUpdateIpFilterRule(ctx context.Context, id IpfilterruleId, input IpFilterRule) (*http.Request, error) {
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

// responderForNamespacesCreateOrUpdateIpFilterRule handles the response to the NamespacesCreateOrUpdateIpFilterRule request. The method always
// closes the http.Response Body.
func (c IpFilterRulesClient) responderForNamespacesCreateOrUpdateIpFilterRule(resp *http.Response) (result NamespacesCreateOrUpdateIpFilterRuleResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
