package ipfilterrules

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type NamespacesDeleteIpFilterRuleResponse struct {
	HttpResponse *http.Response
}

// NamespacesDeleteIpFilterRule ...
func (c IpFilterRulesClient) NamespacesDeleteIpFilterRule(ctx context.Context, id IpfilterruleId) (result NamespacesDeleteIpFilterRuleResponse, err error) {
	req, err := c.preparerForNamespacesDeleteIpFilterRule(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "ipfilterrules.IpFilterRulesClient", "NamespacesDeleteIpFilterRule", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "ipfilterrules.IpFilterRulesClient", "NamespacesDeleteIpFilterRule", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForNamespacesDeleteIpFilterRule(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "ipfilterrules.IpFilterRulesClient", "NamespacesDeleteIpFilterRule", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForNamespacesDeleteIpFilterRule prepares the NamespacesDeleteIpFilterRule request.
func (c IpFilterRulesClient) preparerForNamespacesDeleteIpFilterRule(ctx context.Context, id IpfilterruleId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForNamespacesDeleteIpFilterRule handles the response to the NamespacesDeleteIpFilterRule request. The method always
// closes the http.Response Body.
func (c IpFilterRulesClient) responderForNamespacesDeleteIpFilterRule(resp *http.Response) (result NamespacesDeleteIpFilterRuleResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
