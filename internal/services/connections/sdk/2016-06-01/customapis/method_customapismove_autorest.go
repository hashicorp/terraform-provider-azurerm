package customapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type CustomApisMoveOperationResponse struct {
	HttpResponse *http.Response
}

// CustomApisMove ...
func (c CustomAPIsClient) CustomApisMove(ctx context.Context, id CustomApiId, input CustomApiReference) (result CustomApisMoveOperationResponse, err error) {
	req, err := c.preparerForCustomApisMove(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisMove", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisMove", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCustomApisMove(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisMove", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCustomApisMove prepares the CustomApisMove request.
func (c CustomAPIsClient) preparerForCustomApisMove(ctx context.Context, id CustomApiId, input CustomApiReference) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/move", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCustomApisMove handles the response to the CustomApisMove request. The method always
// closes the http.Response Body.
func (c CustomAPIsClient) responderForCustomApisMove(resp *http.Response) (result CustomApisMoveOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
