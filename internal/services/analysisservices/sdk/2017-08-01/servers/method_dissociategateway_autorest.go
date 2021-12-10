package servers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type DissociateGatewayResponse struct {
	HttpResponse *http.Response
}

// DissociateGateway ...
func (c ServersClient) DissociateGateway(ctx context.Context, id ServerId) (result DissociateGatewayResponse, err error) {
	req, err := c.preparerForDissociateGateway(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servers.ServersClient", "DissociateGateway", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "servers.ServersClient", "DissociateGateway", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDissociateGateway(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servers.ServersClient", "DissociateGateway", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDissociateGateway prepares the DissociateGateway request.
func (c ServersClient) preparerForDissociateGateway(ctx context.Context, id ServerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/dissociateGateway", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDissociateGateway handles the response to the DissociateGateway request. The method always
// closes the http.Response Body.
func (c ServersClient) responderForDissociateGateway(resp *http.Response) (result DissociateGatewayResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
