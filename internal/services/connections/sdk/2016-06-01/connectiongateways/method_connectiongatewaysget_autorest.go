package connectiongateways

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ConnectionGatewaysGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ConnectionGatewayDefinition
}

// ConnectionGatewaysGet ...
func (c ConnectionGatewaysClient) ConnectionGatewaysGet(ctx context.Context, id ConnectionGatewayId) (result ConnectionGatewaysGetOperationResponse, err error) {
	req, err := c.preparerForConnectionGatewaysGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewaysGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewaysGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConnectionGatewaysGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewaysGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConnectionGatewaysGet prepares the ConnectionGatewaysGet request.
func (c ConnectionGatewaysClient) preparerForConnectionGatewaysGet(ctx context.Context, id ConnectionGatewayId) (*http.Request, error) {
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

// responderForConnectionGatewaysGet handles the response to the ConnectionGatewaysGet request. The method always
// closes the http.Response Body.
func (c ConnectionGatewaysClient) responderForConnectionGatewaysGet(resp *http.Response) (result ConnectionGatewaysGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
