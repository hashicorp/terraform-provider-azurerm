package customapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type CustomApisDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// CustomApisDelete ...
func (c CustomAPIsClient) CustomApisDelete(ctx context.Context, id CustomApiId) (result CustomApisDeleteOperationResponse, err error) {
	req, err := c.preparerForCustomApisDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCustomApisDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCustomApisDelete prepares the CustomApisDelete request.
func (c CustomAPIsClient) preparerForCustomApisDelete(ctx context.Context, id CustomApiId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCustomApisDelete handles the response to the CustomApisDelete request. The method always
// closes the http.Response Body.
func (c CustomAPIsClient) responderForCustomApisDelete(resp *http.Response) (result CustomApisDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
