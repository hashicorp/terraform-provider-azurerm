package connectiongateways

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

type ConnectionGatewaysListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *ConnectionGatewayDefinitionCollection
}

// ConnectionGatewaysListByResourceGroup ...
func (c ConnectionGatewaysClient) ConnectionGatewaysListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result ConnectionGatewaysListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForConnectionGatewaysListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewaysListByResourceGroup", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewaysListByResourceGroup", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConnectionGatewaysListByResourceGroup(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectiongateways.ConnectionGatewaysClient", "ConnectionGatewaysListByResourceGroup", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConnectionGatewaysListByResourceGroup prepares the ConnectionGatewaysListByResourceGroup request.
func (c ConnectionGatewaysClient) preparerForConnectionGatewaysListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Web/connectionGateways", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForConnectionGatewaysListByResourceGroup handles the response to the ConnectionGatewaysListByResourceGroup request. The method always
// closes the http.Response Body.
func (c ConnectionGatewaysClient) responderForConnectionGatewaysListByResourceGroup(resp *http.Response) (result ConnectionGatewaysListByResourceGroupOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
