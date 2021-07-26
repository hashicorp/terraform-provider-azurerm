package authorizationrulesnamespaces

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type NamespacesRegenerateKeysResponse struct {
	HttpResponse *http.Response
	Model        *AccessKeys
}

// NamespacesRegenerateKeys ...
func (c AuthorizationRulesNamespacesClient) NamespacesRegenerateKeys(ctx context.Context, id AuthorizationRuleId, input RegenerateAccessKeyParameters) (result NamespacesRegenerateKeysResponse, err error) {
	req, err := c.preparerForNamespacesRegenerateKeys(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesRegenerateKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesRegenerateKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForNamespacesRegenerateKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesRegenerateKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForNamespacesRegenerateKeys prepares the NamespacesRegenerateKeys request.
func (c AuthorizationRulesNamespacesClient) preparerForNamespacesRegenerateKeys(ctx context.Context, id AuthorizationRuleId, input RegenerateAccessKeyParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/regenerateKeys", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForNamespacesRegenerateKeys handles the response to the NamespacesRegenerateKeys request. The method always
// closes the http.Response Body.
func (c AuthorizationRulesNamespacesClient) responderForNamespacesRegenerateKeys(resp *http.Response) (result NamespacesRegenerateKeysResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
