package connectiongateways

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ConnectionGatewaysDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// ConnectionGatewaysDelete ...
func (c ConnectionGatewaysClient) ConnectionGatewaysDelete(ctx context.Context, id ConnectionGatewayId) (result ConnectionGatewaysDeleteOperationResponse, err error) {
	req, err := c.preparerForConnectionGatewaysDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewaysDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewaysDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConnectionGatewaysDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewaysDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConnectionGatewaysDelete prepares the ConnectionGatewaysDelete request.
func (c ConnectionGatewaysClient) preparerForConnectionGatewaysDelete(ctx context.Context, id ConnectionGatewayId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForConnectionGatewaysDelete handles the response to the ConnectionGatewaysDelete request. The method always
// closes the http.Response Body.
func (c ConnectionGatewaysClient) responderForConnectionGatewaysDelete(resp *http.Response) (result ConnectionGatewaysDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
