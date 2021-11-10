package get

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type PrivateLinkResourcesGetResponse struct {
	HttpResponse *http.Response
	Model        *GroupIdInformation
}

// PrivateLinkResourcesGet ...
func (c GETClient) PrivateLinkResourcesGet(ctx context.Context, id PrivateLinkResourceId) (result PrivateLinkResourcesGetResponse, err error) {
	req, err := c.preparerForPrivateLinkResourcesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "get.GETClient", "PrivateLinkResourcesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "get.GETClient", "PrivateLinkResourcesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPrivateLinkResourcesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "get.GETClient", "PrivateLinkResourcesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPrivateLinkResourcesGet prepares the PrivateLinkResourcesGet request.
func (c GETClient) preparerForPrivateLinkResourcesGet(ctx context.Context, id PrivateLinkResourceId) (*http.Request, error) {
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

// responderForPrivateLinkResourcesGet handles the response to the PrivateLinkResourcesGet request. The method always
// closes the http.Response Body.
func (c GETClient) responderForPrivateLinkResourcesGet(resp *http.Response) (result PrivateLinkResourcesGetResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
