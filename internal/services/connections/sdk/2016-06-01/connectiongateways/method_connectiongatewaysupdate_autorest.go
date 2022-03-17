package connectiongateways

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ConnectionGatewaysUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *ConnectionGatewayDefinition
}

// ConnectionGatewaysUpdate ...
func (c ConnectionGatewaysClient) ConnectionGatewaysUpdate(ctx context.Context, id ConnectionGatewayId, input ConnectionGatewayDefinition) (result ConnectionGatewaysUpdateOperationResponse, err error) {
	req, err := c.preparerForConnectionGatewaysUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewaysUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewaysUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConnectionGatewaysUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewaysUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConnectionGatewaysUpdate prepares the ConnectionGatewaysUpdate request.
func (c ConnectionGatewaysClient) preparerForConnectionGatewaysUpdate(ctx context.Context, id ConnectionGatewayId, input ConnectionGatewayDefinition) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForConnectionGatewaysUpdate handles the response to the ConnectionGatewaysUpdate request. The method always
// closes the http.Response Body.
func (c ConnectionGatewaysClient) responderForConnectionGatewaysUpdate(resp *http.Response) (result ConnectionGatewaysUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
