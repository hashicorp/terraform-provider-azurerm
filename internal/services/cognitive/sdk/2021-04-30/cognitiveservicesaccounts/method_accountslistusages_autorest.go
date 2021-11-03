package cognitiveservicesaccounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type AccountsListUsagesResponse struct {
	HttpResponse *http.Response
	Model        *UsageListResult
}

type AccountsListUsagesOptions struct {
	Filter *string
}

func DefaultAccountsListUsagesOptions() AccountsListUsagesOptions {
	return AccountsListUsagesOptions{}
}

func (o AccountsListUsagesOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// AccountsListUsages ...
func (c CognitiveServicesAccountsClient) AccountsListUsages(ctx context.Context, id AccountId, options AccountsListUsagesOptions) (result AccountsListUsagesResponse, err error) {
	req, err := c.preparerForAccountsListUsages(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListUsages", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListUsages", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAccountsListUsages(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListUsages", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAccountsListUsages prepares the AccountsListUsages request.
func (c CognitiveServicesAccountsClient) preparerForAccountsListUsages(ctx context.Context, id AccountId, options AccountsListUsagesOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/usages", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAccountsListUsages handles the response to the AccountsListUsages request. The method always
// closes the http.Response Body.
func (c CognitiveServicesAccountsClient) responderForAccountsListUsages(resp *http.Response) (result AccountsListUsagesResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
