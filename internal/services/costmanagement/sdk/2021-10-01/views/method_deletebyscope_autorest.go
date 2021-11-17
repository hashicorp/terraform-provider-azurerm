package views

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type DeleteByScopeResponse struct {
	HttpResponse *http.Response
}

// DeleteByScope ...
func (c ViewsClient) DeleteByScope(ctx context.Context, id ScopedViewId) (result DeleteByScopeResponse, err error) {
	req, err := c.preparerForDeleteByScope(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "views.ViewsClient", "DeleteByScope", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "views.ViewsClient", "DeleteByScope", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDeleteByScope(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "views.ViewsClient", "DeleteByScope", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDeleteByScope prepares the DeleteByScope request.
func (c ViewsClient) preparerForDeleteByScope(ctx context.Context, id ScopedViewId) (*http.Request, error) {
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

// responderForDeleteByScope handles the response to the DeleteByScope request. The method always
// closes the http.Response Body.
func (c ViewsClient) responderForDeleteByScope(resp *http.Response) (result DeleteByScopeResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
