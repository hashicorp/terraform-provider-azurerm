package keys

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type CreateIfNotExistResponse struct {
	HttpResponse *http.Response
	Model        *Key
}

// CreateIfNotExist ...
func (c KeysClient) CreateIfNotExist(ctx context.Context, id KeyId, input KeyCreateParameters) (result CreateIfNotExistResponse, err error) {
	req, err := c.preparerForCreateIfNotExist(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "keys.KeysClient", "CreateIfNotExist", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "keys.KeysClient", "CreateIfNotExist", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCreateIfNotExist(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "keys.KeysClient", "CreateIfNotExist", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCreateIfNotExist prepares the CreateIfNotExist request.
func (c KeysClient) preparerForCreateIfNotExist(ctx context.Context, id KeyId, input KeyCreateParameters) (*http.Request, error) {
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

// responderForCreateIfNotExist handles the response to the CreateIfNotExist request. The method always
// closes the http.Response Body.
func (c KeysClient) responderForCreateIfNotExist(resp *http.Response) (result CreateIfNotExistResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
