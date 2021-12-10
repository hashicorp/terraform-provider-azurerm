package signalr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type RestartResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Restart ...
func (c SignalRClient) Restart(ctx context.Context, id SignalRId) (result RestartResponse, err error) {
	req, err := c.preparerForRestart(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "Restart", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForRestart(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "Restart", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// RestartThenPoll performs Restart then polls until it's completed
func (c SignalRClient) RestartThenPoll(ctx context.Context, id SignalRId) error {
	result, err := c.Restart(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Restart: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Restart: %+v", err)
	}

	return nil
}

// preparerForRestart prepares the Restart request.
func (c SignalRClient) preparerForRestart(ctx context.Context, id SignalRId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/restart", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForRestart sends the Restart request. The method will close the
// http.Response Body if it receives an error.
func (c SignalRClient) senderForRestart(ctx context.Context, req *http.Request) (future RestartResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
