package cognitiveservicesaccounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type AccountsListKeysResponse struct {
	HttpResponse *http.Response
	Model        *ApiKeys
}

// AccountsListKeys ...
func (c CognitiveServicesAccountsClient) AccountsListKeys(ctx context.Context, id AccountId) (result AccountsListKeysResponse, err error) {
	req, err := c.preparerForAccountsListKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAccountsListKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAccountsListKeys prepares the AccountsListKeys request.
func (c CognitiveServicesAccountsClient) preparerForAccountsListKeys(ctx context.Context, id AccountId) (*http.Request, error) {
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

// responderForAccountsListKeys handles the response to the AccountsListKeys request. The method always
// closes the http.Response Body.
func (c CognitiveServicesAccountsClient) responderForAccountsListKeys(resp *http.Response) (result AccountsListKeysResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
