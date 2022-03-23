package connectiongateways

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ConnectionGatewayInstallationsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *ConnectionGatewayInstallationDefinitionCollection
}

// ConnectionGatewayInstallationsList ...
func (c ConnectionGatewaysClient) ConnectionGatewayInstallationsList(ctx context.Context, id LocationId) (result ConnectionGatewayInstallationsListOperationResponse, err error) {
	req, err := c.preparerForConnectionGatewayInstallationsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewayInstallationsList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewayInstallationsList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConnectionGatewayInstallationsList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewayInstallationsList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConnectionGatewayInstallationsList prepares the ConnectionGatewayInstallationsList request.
func (c ConnectionGatewaysClient) preparerForConnectionGatewayInstallationsList(ctx context.Context, id LocationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/connectionGatewayInstallations", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForConnectionGatewayInstallationsList handles the response to the ConnectionGatewayInstallationsList request. The method always
// closes the http.Response Body.
func (c ConnectionGatewaysClient) responderForConnectionGatewayInstallationsList(resp *http.Response) (result ConnectionGatewayInstallationsListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
