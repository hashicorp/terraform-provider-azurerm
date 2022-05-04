package networkrulesets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type NamespacesListNetworkRuleSetOperationResponse struct {
	HttpResponse *http.Response
	Model        *NetworkRuleSetListResult
}

// NamespacesListNetworkRuleSet ...
func (c NetworkRuleSetsClient) NamespacesListNetworkRuleSet(ctx context.Context, id NamespaceId) (result NamespacesListNetworkRuleSetOperationResponse, err error) {
	req, err := c.preparerForNamespacesListNetworkRuleSet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "networkrulesets.NetworkRuleSetsClient", "NamespacesListNetworkRuleSet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "networkrulesets.NetworkRuleSetsClient", "NamespacesListNetworkRuleSet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForNamespacesListNetworkRuleSet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "networkrulesets.NetworkRuleSetsClient", "NamespacesListNetworkRuleSet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForNamespacesListNetworkRuleSet prepares the NamespacesListNetworkRuleSet request.
func (c NetworkRuleSetsClient) preparerForNamespacesListNetworkRuleSet(ctx context.Context, id NamespaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/networkRuleSets", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForNamespacesListNetworkRuleSet handles the response to the NamespacesListNetworkRuleSet request. The method always
// closes the http.Response Body.
func (c NetworkRuleSetsClient) responderForNamespacesListNetworkRuleSet(resp *http.Response) (result NamespacesListNetworkRuleSetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
