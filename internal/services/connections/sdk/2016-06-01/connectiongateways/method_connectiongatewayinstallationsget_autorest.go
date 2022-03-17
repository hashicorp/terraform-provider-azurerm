package connectiongateways

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ConnectionGatewayInstallationsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ConnectionGatewayInstallationDefinition
}

// ConnectionGatewayInstallationsGet ...
func (c ConnectionGatewaysClient) ConnectionGatewayInstallationsGet(ctx context.Context, id ConnectionGatewayInstallationId) (result ConnectionGatewayInstallationsGetOperationResponse, err error) {
	req, err := c.preparerForConnectionGatewayInstallationsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewayInstallationsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewayInstallationsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConnectionGatewayInstallationsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewayInstallationsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConnectionGatewayInstallationsGet prepares the ConnectionGatewayInstallationsGet request.
func (c ConnectionGatewaysClient) preparerForConnectionGatewayInstallationsGet(ctx context.Context, id ConnectionGatewayInstallationId) (*http.Request, error) {
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

// responderForConnectionGatewayInstallationsGet handles the response to the ConnectionGatewayInstallationsGet request. The method always
// closes the http.Response Body.
func (c ConnectionGatewaysClient) responderForConnectionGatewayInstallationsGet(resp *http.Response) (result ConnectionGatewayInstallationsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
