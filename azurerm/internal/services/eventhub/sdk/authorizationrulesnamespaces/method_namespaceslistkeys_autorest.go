package authorizationrulesnamespaces

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type NamespacesListKeysResponse struct {
	HttpResponse *http.Response
	Model        *AccessKeys
}

// NamespacesListKeys ...
func (c AuthorizationRulesNamespacesClient) NamespacesListKeys(ctx context.Context, id AuthorizationRuleId) (result NamespacesListKeysResponse, err error) {
	req, err := c.preparerForNamespacesListKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesListKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesListKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForNamespacesListKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesListKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForNamespacesListKeys prepares the NamespacesListKeys request.
func (c AuthorizationRulesNamespacesClient) preparerForNamespacesListKeys(ctx context.Context, id AuthorizationRuleId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForNamespacesListKeys handles the response to the NamespacesListKeys request. The method always
// closes the http.Response Body.
func (c AuthorizationRulesNamespacesClient) responderForNamespacesListKeys(resp *http.Response) (result NamespacesListKeysResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
