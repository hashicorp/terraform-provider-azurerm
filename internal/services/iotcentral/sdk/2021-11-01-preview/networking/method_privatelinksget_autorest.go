package networking

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type PrivateLinksGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *PrivateLinkResource
}

// PrivateLinksGet ...
func (c NetworkingClient) PrivateLinksGet(ctx context.Context, id PrivateLinkResourceId) (result PrivateLinksGetOperationResponse, err error) {
	req, err := c.preparerForPrivateLinksGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "networking.NetworkingClient", "PrivateLinksGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "networking.NetworkingClient", "PrivateLinksGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPrivateLinksGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "networking.NetworkingClient", "PrivateLinksGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPrivateLinksGet prepares the PrivateLinksGet request.
func (c NetworkingClient) preparerForPrivateLinksGet(ctx context.Context, id PrivateLinkResourceId) (*http.Request, error) {
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

// responderForPrivateLinksGet handles the response to the PrivateLinksGet request. The method always
// closes the http.Response Body.
func (c NetworkingClient) responderForPrivateLinksGet(resp *http.Response) (result PrivateLinksGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
