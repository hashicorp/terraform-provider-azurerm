package kubernetes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type ConnectedClusterDeleteResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ConnectedClusterDelete ...
func (c KubernetesClient) ConnectedClusterDelete(ctx context.Context, id ConnectedClusterId) (result ConnectedClusterDeleteResponse, err error) {
	req, err := c.preparerForConnectedClusterDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "kubernetes.KubernetesClient", "ConnectedClusterDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForConnectedClusterDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "kubernetes.KubernetesClient", "ConnectedClusterDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ConnectedClusterDeleteThenPoll performs ConnectedClusterDelete then polls until it's completed
func (c KubernetesClient) ConnectedClusterDeleteThenPoll(ctx context.Context, id ConnectedClusterId) error {
	result, err := c.ConnectedClusterDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ConnectedClusterDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ConnectedClusterDelete: %+v", err)
	}

	return nil
}

// preparerForConnectedClusterDelete prepares the ConnectedClusterDelete request.
func (c KubernetesClient) preparerForConnectedClusterDelete(ctx context.Context, id ConnectedClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForConnectedClusterDelete sends the ConnectedClusterDelete request. The method will close the
// http.Response Body if it receives an error.
func (c KubernetesClient) senderForConnectedClusterDelete(ctx context.Context, req *http.Request) (future ConnectedClusterDeleteResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
