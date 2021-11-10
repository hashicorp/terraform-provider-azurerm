package networkrulesets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type NamespacesCreateOrUpdateNetworkRuleSetResponse struct {
	HttpResponse *http.Response
	Model        *NetworkRuleSet
}

// NamespacesCreateOrUpdateNetworkRuleSet ...
func (c NetworkRuleSetsClient) NamespacesCreateOrUpdateNetworkRuleSet(ctx context.Context, id NamespaceId, input NetworkRuleSet) (result NamespacesCreateOrUpdateNetworkRuleSetResponse, err error) {
	req, err := c.preparerForNamespacesCreateOrUpdateNetworkRuleSet(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "networkrulesets.NetworkRuleSetsClient", "NamespacesCreateOrUpdateNetworkRuleSet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "networkrulesets.NetworkRuleSetsClient", "NamespacesCreateOrUpdateNetworkRuleSet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForNamespacesCreateOrUpdateNetworkRuleSet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "networkrulesets.NetworkRuleSetsClient", "NamespacesCreateOrUpdateNetworkRuleSet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForNamespacesCreateOrUpdateNetworkRuleSet prepares the NamespacesCreateOrUpdateNetworkRuleSet request.
func (c NetworkRuleSetsClient) preparerForNamespacesCreateOrUpdateNetworkRuleSet(ctx context.Context, id NamespaceId, input NetworkRuleSet) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/networkRuleSets/default", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForNamespacesCreateOrUpdateNetworkRuleSet handles the response to the NamespacesCreateOrUpdateNetworkRuleSet request. The method always
// closes the http.Response Body.
func (c NetworkRuleSetsClient) responderForNamespacesCreateOrUpdateNetworkRuleSet(resp *http.Response) (result NamespacesCreateOrUpdateNetworkRuleSetResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
