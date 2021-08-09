package signalr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type UpdateResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Update ...
func (c SignalRClient) Update(ctx context.Context, id SignalRId, input SignalRResource) (result UpdateResponse, err error) {
	req, err := c.preparerForUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "Update", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "Update", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// UpdateThenPoll performs Update then polls until it's completed
func (c SignalRClient) UpdateThenPoll(ctx context.Context, id SignalRId, input SignalRResource) error {
	result, err := c.Update(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Update: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Update: %+v", err)
	}

	return nil
}

// preparerForUpdate prepares the Update request.
func (c SignalRClient) preparerForUpdate(ctx context.Context, id SignalRId, input SignalRResource) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForUpdate sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (c SignalRClient) senderForUpdate(ctx context.Context, req *http.Request) (future UpdateResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
