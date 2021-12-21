package datalakestoreaccounts

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type AddResponse struct {
	HttpResponse *http.Response
}

// Add ...
func (c DataLakeStoreAccountsClient) Add(ctx context.Context, id DataLakeStoreAccountId, input AddDataLakeStoreParameters) (result AddResponse, err error) {
	req, err := c.preparerForAdd(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datalakestoreaccounts.DataLakeStoreAccountsClient", "Add", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "datalakestoreaccounts.DataLakeStoreAccountsClient", "Add", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAdd(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datalakestoreaccounts.DataLakeStoreAccountsClient", "Add", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAdd prepares the Add request.
func (c DataLakeStoreAccountsClient) preparerForAdd(ctx context.Context, id DataLakeStoreAccountId, input AddDataLakeStoreParameters) (*http.Request, error) {
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

// responderForAdd handles the response to the Add request. The method always
// closes the http.Response Body.
func (c DataLakeStoreAccountsClient) responderForAdd(resp *http.Response) (result AddResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
