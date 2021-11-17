package views

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type CreateOrUpdateByScopeResponse struct {
	HttpResponse *http.Response
	Model        *View
}

// CreateOrUpdateByScope ...
func (c ViewsClient) CreateOrUpdateByScope(ctx context.Context, id ScopedViewId, input View) (result CreateOrUpdateByScopeResponse, err error) {
	req, err := c.preparerForCreateOrUpdateByScope(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "views.ViewsClient", "CreateOrUpdateByScope", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "views.ViewsClient", "CreateOrUpdateByScope", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCreateOrUpdateByScope(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "views.ViewsClient", "CreateOrUpdateByScope", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCreateOrUpdateByScope prepares the CreateOrUpdateByScope request.
func (c ViewsClient) preparerForCreateOrUpdateByScope(ctx context.Context, id ScopedViewId, input View) (*http.Request, error) {
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

// responderForCreateOrUpdateByScope handles the response to the CreateOrUpdateByScope request. The method always
// closes the http.Response Body.
func (c ViewsClient) responderForCreateOrUpdateByScope(resp *http.Response) (result CreateOrUpdateByScopeResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
