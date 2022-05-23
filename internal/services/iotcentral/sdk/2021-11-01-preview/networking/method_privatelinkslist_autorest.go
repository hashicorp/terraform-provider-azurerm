package networking

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type PrivateLinksListOperationResponse struct {
	HttpResponse *http.Response
	Model        *PrivateLinkResourceListResult
}

// PrivateLinksList ...
func (c NetworkingClient) PrivateLinksList(ctx context.Context, id IotAppId) (result PrivateLinksListOperationResponse, err error) {
	req, err := c.preparerForPrivateLinksList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "networking.NetworkingClient", "PrivateLinksList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "networking.NetworkingClient", "PrivateLinksList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPrivateLinksList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "networking.NetworkingClient", "PrivateLinksList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPrivateLinksList prepares the PrivateLinksList request.
func (c NetworkingClient) preparerForPrivateLinksList(ctx context.Context, id IotAppId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/privateLinkResources", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForPrivateLinksList handles the response to the PrivateLinksList request. The method always
// closes the http.Response Body.
func (c NetworkingClient) responderForPrivateLinksList(resp *http.Response) (result PrivateLinksListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
