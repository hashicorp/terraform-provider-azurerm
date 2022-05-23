package networking

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type PrivateEndpointConnectionsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *PrivateEndpointConnectionListResult
}

// PrivateEndpointConnectionsList ...
func (c NetworkingClient) PrivateEndpointConnectionsList(ctx context.Context, id IotAppId) (result PrivateEndpointConnectionsListOperationResponse, err error) {
	req, err := c.preparerForPrivateEndpointConnectionsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "networking.NetworkingClient", "PrivateEndpointConnectionsList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "networking.NetworkingClient", "PrivateEndpointConnectionsList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPrivateEndpointConnectionsList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "networking.NetworkingClient", "PrivateEndpointConnectionsList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPrivateEndpointConnectionsList prepares the PrivateEndpointConnectionsList request.
func (c NetworkingClient) preparerForPrivateEndpointConnectionsList(ctx context.Context, id IotAppId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/privateEndpointConnections", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForPrivateEndpointConnectionsList handles the response to the PrivateEndpointConnectionsList request. The method always
// closes the http.Response Body.
func (c NetworkingClient) responderForPrivateEndpointConnectionsList(resp *http.Response) (result PrivateEndpointConnectionsListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
