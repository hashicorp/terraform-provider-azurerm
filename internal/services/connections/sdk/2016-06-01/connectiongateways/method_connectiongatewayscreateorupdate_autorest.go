package connectiongateways

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ConnectionGatewaysCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *ConnectionGatewayDefinition
}

// ConnectionGatewaysCreateOrUpdate ...
func (c ConnectionGatewaysClient) ConnectionGatewaysCreateOrUpdate(ctx context.Context, id ConnectionGatewayId, input ConnectionGatewayDefinition) (result ConnectionGatewaysCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForConnectionGatewaysCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewaysCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewaysCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConnectionGatewaysCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewaysCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConnectionGatewaysCreateOrUpdate prepares the ConnectionGatewaysCreateOrUpdate request.
func (c ConnectionGatewaysClient) preparerForConnectionGatewaysCreateOrUpdate(ctx context.Context, id ConnectionGatewayId, input ConnectionGatewayDefinition) (*http.Request, error) {
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

// responderForConnectionGatewaysCreateOrUpdate handles the response to the ConnectionGatewaysCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c ConnectionGatewaysClient) responderForConnectionGatewaysCreateOrUpdate(resp *http.Response) (result ConnectionGatewaysCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
