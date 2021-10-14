package networkrulesets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type NamespacesGetNetworkRuleSetResponse struct {
	HttpResponse *http.Response
	Model        *NetworkRuleSet
}

// NamespacesGetNetworkRuleSet ...
func (c NetworkRuleSetsClient) NamespacesGetNetworkRuleSet(ctx context.Context, id NamespaceId) (result NamespacesGetNetworkRuleSetResponse, err error) {
	req, err := c.preparerForNamespacesGetNetworkRuleSet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "networkrulesets.NetworkRuleSetsClient", "NamespacesGetNetworkRuleSet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "networkrulesets.NetworkRuleSetsClient", "NamespacesGetNetworkRuleSet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForNamespacesGetNetworkRuleSet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "networkrulesets.NetworkRuleSetsClient", "NamespacesGetNetworkRuleSet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForNamespacesGetNetworkRuleSet prepares the NamespacesGetNetworkRuleSet request.
func (c NetworkRuleSetsClient) preparerForNamespacesGetNetworkRuleSet(ctx context.Context, id NamespaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/networkRuleSets/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForNamespacesGetNetworkRuleSet handles the response to the NamespacesGetNetworkRuleSet request. The method always
// closes the http.Response Body.
func (c NetworkRuleSetsClient) responderForNamespacesGetNetworkRuleSet(resp *http.Response) (result NamespacesGetNetworkRuleSetResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
