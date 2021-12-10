package signalr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type PrivateEndpointConnectionsDeleteResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// PrivateEndpointConnectionsDelete ...
func (c SignalRClient) PrivateEndpointConnectionsDelete(ctx context.Context, id PrivateEndpointConnectionId) (result PrivateEndpointConnectionsDeleteResponse, err error) {
	req, err := c.preparerForPrivateEndpointConnectionsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "PrivateEndpointConnectionsDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPrivateEndpointConnectionsDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "PrivateEndpointConnectionsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PrivateEndpointConnectionsDeleteThenPoll performs PrivateEndpointConnectionsDelete then polls until it's completed
func (c SignalRClient) PrivateEndpointConnectionsDeleteThenPoll(ctx context.Context, id PrivateEndpointConnectionId) error {
	result, err := c.PrivateEndpointConnectionsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing PrivateEndpointConnectionsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after PrivateEndpointConnectionsDelete: %+v", err)
	}

	return nil
}

// preparerForPrivateEndpointConnectionsDelete prepares the PrivateEndpointConnectionsDelete request.
func (c SignalRClient) preparerForPrivateEndpointConnectionsDelete(ctx context.Context, id PrivateEndpointConnectionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForPrivateEndpointConnectionsDelete sends the PrivateEndpointConnectionsDelete request. The method will close the
// http.Response Body if it receives an error.
func (c SignalRClient) senderForPrivateEndpointConnectionsDelete(ctx context.Context, req *http.Request) (future PrivateEndpointConnectionsDeleteResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
