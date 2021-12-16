package redisenterprise

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type DatabasesGetResponse struct {
	HttpResponse *http.Response
	Model        *Database
}

// DatabasesGet ...
func (c RedisEnterpriseClient) DatabasesGet(ctx context.Context, id DatabaseId) (result DatabasesGetResponse, err error) {
	req, err := c.preparerForDatabasesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDatabasesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDatabasesGet prepares the DatabasesGet request.
func (c RedisEnterpriseClient) preparerForDatabasesGet(ctx context.Context, id DatabaseId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDatabasesGet handles the response to the DatabasesGet request. The method always
// closes the http.Response Body.
func (c RedisEnterpriseClient) responderForDatabasesGet(resp *http.Response) (result DatabasesGetResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
