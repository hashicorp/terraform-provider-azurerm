package views

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type GetByScopeResponse struct {
	HttpResponse *http.Response
	Model        *View
}

// GetByScope ...
func (c ViewsClient) GetByScope(ctx context.Context, id ScopedViewId) (result GetByScopeResponse, err error) {
	req, err := c.preparerForGetByScope(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "views.ViewsClient", "GetByScope", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "views.ViewsClient", "GetByScope", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetByScope(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "views.ViewsClient", "GetByScope", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetByScope prepares the GetByScope request.
func (c ViewsClient) preparerForGetByScope(ctx context.Context, id ScopedViewId) (*http.Request, error) {
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

// responderForGetByScope handles the response to the GetByScope request. The method always
// closes the http.Response Body.
func (c ViewsClient) responderForGetByScope(resp *http.Response) (result GetByScopeResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
