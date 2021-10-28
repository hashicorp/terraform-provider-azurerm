package mhsmprivateendpointconnections

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type PutResponse struct {
	HttpResponse *http.Response
	Model        *MHSMPrivateEndpointConnection
}

// Put ...
func (c MHSMPrivateEndpointConnectionsClient) Put(ctx context.Context, id PrivateEndpointConnectionId, input MHSMPrivateEndpointConnection) (result PutResponse, err error) {
	req, err := c.preparerForPut(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mhsmprivateendpointconnections.MHSMPrivateEndpointConnectionsClient", "Put", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "mhsmprivateendpointconnections.MHSMPrivateEndpointConnectionsClient", "Put", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPut(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mhsmprivateendpointconnections.MHSMPrivateEndpointConnectionsClient", "Put", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPut prepares the Put request.
func (c MHSMPrivateEndpointConnectionsClient) preparerForPut(ctx context.Context, id PrivateEndpointConnectionId, input MHSMPrivateEndpointConnection) (*http.Request, error) {
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

// responderForPut handles the response to the Put request. The method always
// closes the http.Response Body.
func (c MHSMPrivateEndpointConnectionsClient) responderForPut(resp *http.Response) (result PutResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
