package get

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type OutboundNetworkDependenciesEndpointsListResponse struct {
	HttpResponse *http.Response
	Model        *[]OutboundEnvironmentEndpoint
}

// OutboundNetworkDependenciesEndpointsList ...
func (c GETClient) OutboundNetworkDependenciesEndpointsList(ctx context.Context, id WorkspaceId) (result OutboundNetworkDependenciesEndpointsListResponse, err error) {
	req, err := c.preparerForOutboundNetworkDependenciesEndpointsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "get.GETClient", "OutboundNetworkDependenciesEndpointsList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "get.GETClient", "OutboundNetworkDependenciesEndpointsList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForOutboundNetworkDependenciesEndpointsList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "get.GETClient", "OutboundNetworkDependenciesEndpointsList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForOutboundNetworkDependenciesEndpointsList prepares the OutboundNetworkDependenciesEndpointsList request.
func (c GETClient) preparerForOutboundNetworkDependenciesEndpointsList(ctx context.Context, id WorkspaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/outboundNetworkDependenciesEndpoints", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForOutboundNetworkDependenciesEndpointsList handles the response to the OutboundNetworkDependenciesEndpointsList request. The method always
// closes the http.Response Body.
func (c GETClient) responderForOutboundNetworkDependenciesEndpointsList(resp *http.Response) (result OutboundNetworkDependenciesEndpointsListResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
