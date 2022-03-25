package customapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type CustomApisGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *CustomApiDefinition
}

// CustomApisGet ...
func (c CustomAPIsClient) CustomApisGet(ctx context.Context, id CustomApiId) (result CustomApisGetOperationResponse, err error) {
	req, err := c.preparerForCustomApisGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCustomApisGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCustomApisGet prepares the CustomApisGet request.
func (c CustomAPIsClient) preparerForCustomApisGet(ctx context.Context, id CustomApiId) (*http.Request, error) {
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

// responderForCustomApisGet handles the response to the CustomApisGet request. The method always
// closes the http.Response Body.
func (c CustomAPIsClient) responderForCustomApisGet(resp *http.Response) (result CustomApisGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
