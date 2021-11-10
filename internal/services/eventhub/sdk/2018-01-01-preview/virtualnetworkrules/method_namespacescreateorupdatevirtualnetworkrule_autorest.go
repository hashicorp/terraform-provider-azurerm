package virtualnetworkrules

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type NamespacesCreateOrUpdateVirtualNetworkRuleResponse struct {
	HttpResponse *http.Response
	Model        *VirtualNetworkRule
}

// NamespacesCreateOrUpdateVirtualNetworkRule ...
func (c VirtualNetworkRulesClient) NamespacesCreateOrUpdateVirtualNetworkRule(ctx context.Context, id VirtualnetworkruleId, input VirtualNetworkRule) (result NamespacesCreateOrUpdateVirtualNetworkRuleResponse, err error) {
	req, err := c.preparerForNamespacesCreateOrUpdateVirtualNetworkRule(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualnetworkrules.VirtualNetworkRulesClient", "NamespacesCreateOrUpdateVirtualNetworkRule", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualnetworkrules.VirtualNetworkRulesClient", "NamespacesCreateOrUpdateVirtualNetworkRule", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForNamespacesCreateOrUpdateVirtualNetworkRule(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualnetworkrules.VirtualNetworkRulesClient", "NamespacesCreateOrUpdateVirtualNetworkRule", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForNamespacesCreateOrUpdateVirtualNetworkRule prepares the NamespacesCreateOrUpdateVirtualNetworkRule request.
func (c VirtualNetworkRulesClient) preparerForNamespacesCreateOrUpdateVirtualNetworkRule(ctx context.Context, id VirtualnetworkruleId, input VirtualNetworkRule) (*http.Request, error) {
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

// responderForNamespacesCreateOrUpdateVirtualNetworkRule handles the response to the NamespacesCreateOrUpdateVirtualNetworkRule request. The method always
// closes the http.Response Body.
func (c VirtualNetworkRulesClient) responderForNamespacesCreateOrUpdateVirtualNetworkRule(resp *http.Response) (result NamespacesCreateOrUpdateVirtualNetworkRuleResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
