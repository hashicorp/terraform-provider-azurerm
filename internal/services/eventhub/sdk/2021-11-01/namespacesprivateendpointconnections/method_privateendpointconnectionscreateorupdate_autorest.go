package namespacesprivateendpointconnections

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type PrivateEndpointConnectionsCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *PrivateEndpointConnection
}

// PrivateEndpointConnectionsCreateOrUpdate ...
func (c NamespacesPrivateEndpointConnectionsClient) PrivateEndpointConnectionsCreateOrUpdate(ctx context.Context, id PrivateEndpointConnectionId, input PrivateEndpointConnection) (result PrivateEndpointConnectionsCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForPrivateEndpointConnectionsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespacesprivateendpointconnections.NamespacesPrivateEndpointConnectionsClient", "PrivateEndpointConnectionsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespacesprivateendpointconnections.NamespacesPrivateEndpointConnectionsClient", "PrivateEndpointConnectionsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPrivateEndpointConnectionsCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespacesprivateendpointconnections.NamespacesPrivateEndpointConnectionsClient", "PrivateEndpointConnectionsCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPrivateEndpointConnectionsCreateOrUpdate prepares the PrivateEndpointConnectionsCreateOrUpdate request.
func (c NamespacesPrivateEndpointConnectionsClient) preparerForPrivateEndpointConnectionsCreateOrUpdate(ctx context.Context, id PrivateEndpointConnectionId, input PrivateEndpointConnection) (*http.Request, error) {
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

// responderForPrivateEndpointConnectionsCreateOrUpdate handles the response to the PrivateEndpointConnectionsCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c NamespacesPrivateEndpointConnectionsClient) responderForPrivateEndpointConnectionsCreateOrUpdate(resp *http.Response) (result PrivateEndpointConnectionsCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
