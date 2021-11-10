package virtualnetworkrules

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type NamespacesGetVirtualNetworkRuleResponse struct {
	HttpResponse *http.Response
	Model        *VirtualNetworkRule
}

// NamespacesGetVirtualNetworkRule ...
func (c VirtualNetworkRulesClient) NamespacesGetVirtualNetworkRule(ctx context.Context, id VirtualnetworkruleId) (result NamespacesGetVirtualNetworkRuleResponse, err error) {
	req, err := c.preparerForNamespacesGetVirtualNetworkRule(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualnetworkrules.VirtualNetworkRulesClient", "NamespacesGetVirtualNetworkRule", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualnetworkrules.VirtualNetworkRulesClient", "NamespacesGetVirtualNetworkRule", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForNamespacesGetVirtualNetworkRule(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualnetworkrules.VirtualNetworkRulesClient", "NamespacesGetVirtualNetworkRule", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForNamespacesGetVirtualNetworkRule prepares the NamespacesGetVirtualNetworkRule request.
func (c VirtualNetworkRulesClient) preparerForNamespacesGetVirtualNetworkRule(ctx context.Context, id VirtualnetworkruleId) (*http.Request, error) {
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

// responderForNamespacesGetVirtualNetworkRule handles the response to the NamespacesGetVirtualNetworkRule request. The method always
// closes the http.Response Body.
func (c VirtualNetworkRulesClient) responderForNamespacesGetVirtualNetworkRule(resp *http.Response) (result NamespacesGetVirtualNetworkRuleResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
