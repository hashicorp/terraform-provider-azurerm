package put

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type PrivateEndpointConnectionsCreateResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// PrivateEndpointConnectionsCreate ...
func (c PUTClient) PrivateEndpointConnectionsCreate(ctx context.Context, id PrivateEndpointConnectionId, input PrivateEndpointConnection) (result PrivateEndpointConnectionsCreateResponse, err error) {
	req, err := c.preparerForPrivateEndpointConnectionsCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "put.PUTClient", "PrivateEndpointConnectionsCreate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPrivateEndpointConnectionsCreate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "put.PUTClient", "PrivateEndpointConnectionsCreate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PrivateEndpointConnectionsCreateThenPoll performs PrivateEndpointConnectionsCreate then polls until it's completed
func (c PUTClient) PrivateEndpointConnectionsCreateThenPoll(ctx context.Context, id PrivateEndpointConnectionId, input PrivateEndpointConnection) error {
	result, err := c.PrivateEndpointConnectionsCreate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing PrivateEndpointConnectionsCreate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after PrivateEndpointConnectionsCreate: %+v", err)
	}

	return nil
}

// preparerForPrivateEndpointConnectionsCreate prepares the PrivateEndpointConnectionsCreate request.
func (c PUTClient) preparerForPrivateEndpointConnectionsCreate(ctx context.Context, id PrivateEndpointConnectionId, input PrivateEndpointConnection) (*http.Request, error) {
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

// senderForPrivateEndpointConnectionsCreate sends the PrivateEndpointConnectionsCreate request. The method will close the
// http.Response Body if it receives an error.
func (c PUTClient) senderForPrivateEndpointConnectionsCreate(ctx context.Context, req *http.Request) (future PrivateEndpointConnectionsCreateResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
