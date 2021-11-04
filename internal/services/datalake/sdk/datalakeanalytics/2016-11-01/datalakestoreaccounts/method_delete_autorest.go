package datalakestoreaccounts

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type DeleteResponse struct {
	HttpResponse *http.Response
}

// Delete ...
func (c DataLakeStoreAccountsClient) Delete(ctx context.Context, id DataLakeStoreAccountId) (result DeleteResponse, err error) {
	req, err := c.preparerForDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datalakestoreaccounts.DataLakeStoreAccountsClient", "Delete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "datalakestoreaccounts.DataLakeStoreAccountsClient", "Delete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datalakestoreaccounts.DataLakeStoreAccountsClient", "Delete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDelete prepares the Delete request.
func (c DataLakeStoreAccountsClient) preparerForDelete(ctx context.Context, id DataLakeStoreAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDelete handles the response to the Delete request. The method always
// closes the http.Response Body.
func (c DataLakeStoreAccountsClient) responderForDelete(resp *http.Response) (result DeleteResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
